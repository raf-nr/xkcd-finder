CREATE TABLE comic_keyword_mapping (
  comic_id INTEGER REFERENCES comic(id) ON DELETE CASCADE,
  keyword_id INTEGER REFERENCES keyword(id) ON DELETE CASCADE,
  PRIMARY KEY (comic_id, keyword_id)
);