package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/raf-nr/xkcd-finder/internal/domain/comic"
)

type ComicRepository struct {
	db *sql.DB
}

func NewComicRepository(db *sql.DB) *ComicRepository {
	return &ComicRepository{
		db: db,
	}
}

func (r *ComicRepository) GetComicsIDs(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM comic")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *ComicRepository) Save(ctx context.Context, c comic.Comic) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = r.saveComic(ctx, tx, c)
	if err != nil {
		return err
	}

	err = r.saveKeywordsAndRelations(ctx, tx, c)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ComicRepository) saveComic(ctx context.Context, tx *sql.Tx, c comic.Comic) error {
	_, err := tx.ExecContext(ctx, `
			INSERT INTO comic (id, image_url) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET image_url = EXCLUDED.image_url
	`, c.ID, c.ImageURL)
	return err
}

func (r *ComicRepository) saveKeywordsAndRelations(ctx context.Context, tx *sql.Tx, c comic.Comic) error {
	keywordsSet := make(map[string]struct{})
	for _, kw := range c.Keywords {
		kw = strings.TrimSpace(kw)
		if kw != "" {
			keywordsSet[kw] = struct{}{}
		}
	}
	if len(keywordsSet) == 0 {
		return nil
	}

	if err := r.insertKeywords(ctx, tx, keywordsSet); err != nil {
		return err
	}

	keywordIDs, err := r.selectKeywordIDs(ctx, tx, keywordsSet)
	if err != nil {
		return err
	}

	return r.insertComicKeywords(ctx, tx, c.ID, keywordIDs)
}

func (r *ComicRepository) insertKeywords(ctx context.Context, tx *sql.Tx, keywordsSet map[string]struct{}) error {
	if len(keywordsSet) == 0 {
		return nil
	}
	var (
		insertValues []string
		insertArgs   []interface{}
		argPos       = 1
	)
	for kw := range keywordsSet {
		insertValues = append(insertValues, "($"+strconv.Itoa(argPos)+")")
		insertArgs = append(insertArgs, kw)
		argPos++
	}

	query := `INSERT INTO keyword (keyword) VALUES ` + strings.Join(insertValues, ",") + `
			ON CONFLICT (keyword) DO NOTHING
	`
	_, err := tx.ExecContext(ctx, query, insertArgs...)
	return err
}

func (r *ComicRepository) selectKeywordIDs(ctx context.Context, tx *sql.Tx, keywordsSet map[string]struct{}) ([]int, error) {
	var (
		selectArgs   []interface{}
		placeholders []string
		i            = 1
	)
	for kw := range keywordsSet {
		selectArgs = append(selectArgs, kw)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}

	query := `SELECT id FROM keyword WHERE keyword IN (` + strings.Join(placeholders, ",") + `)`
	rows, err := tx.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *ComicRepository) insertComicKeywords(ctx context.Context, tx *sql.Tx, comicID int, keywordIDs []int) error {
	if len(keywordIDs) == 0 {
		return nil
	}

	var (
		insertValues []string
		insertArgs   []interface{}
	)
	insertArgs = append(insertArgs, comicID)

	for i, kwID := range keywordIDs {
		insertValues = append(insertValues, "($1,$"+strconv.Itoa(i+2)+")")
		insertArgs = append(insertArgs, kwID)
	}

	query := `INSERT INTO comic_keyword_mapping (comic_id, keyword_id) VALUES ` + strings.Join(insertValues, ",") + `
			ON CONFLICT DO NOTHING
	`
	_, err := tx.ExecContext(ctx, query, insertArgs...)
	return err
}

func (r *ComicRepository) FindTopComicsByKeywords(ctx context.Context, stemmed []string, limit uint) ([]comic.Comic, error) {
	if len(stemmed) == 0 {
		return nil, nil
	}

	query, args := buildTopComicsQuery(stemmed, limit)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanComics(rows)
}

func buildTopComicsQuery(keywords []string, limit uint) (string, []any) {
	var (
		placeholders = make([]string, len(keywords))
		args         = make([]any, len(keywords)+1)
	)

	for i, word := range keywords {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = word
	}

	args[len(keywords)] = limit
	limitPlaceholder := fmt.Sprintf("$%d", len(keywords)+1)

	query := fmt.Sprintf(`
		SELECT c.id, c.image_url, array_agg(k.keyword) AS keywords
		FROM comic c
		JOIN comic_keyword_mapping ckm ON c.id = ckm.comic_id
		JOIN keyword k ON ckm.keyword_id = k.id
		WHERE k.keyword IN (%s)
		GROUP BY c.id
		ORDER BY COUNT(*) DESC
		LIMIT %s
	`, strings.Join(placeholders, ","), limitPlaceholder)

	return query, args
}

func scanComics(rows *sql.Rows) ([]comic.Comic, error) {
	var results []comic.Comic
	for rows.Next() {
		var c comic.Comic
		var keywords []sql.NullString

		if err := rows.Scan(&c.ID, &c.ImageURL, pq.Array(&keywords)); err != nil {
			return nil, err
		}

		for _, kw := range keywords {
			if kw.Valid {
				c.Keywords = append(c.Keywords, kw.String)
			}
		}
		results = append(results, c)
	}
	return results, rows.Err()
}
