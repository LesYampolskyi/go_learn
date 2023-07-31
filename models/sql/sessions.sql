CREATE TABLE session (
  id SERIAL PRIMARY KEY,
  user_id INT UNIQUE,
  token_hash TEXT NOT NULL
);