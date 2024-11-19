DROP DATABASE IF EXISTS excalidb;
CREATE DATABASE excalidb;
-- switch to excalidb database in psql
\c excalidb;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name text NOT NULL,
  username text NOT NULL UNIQUE,
  email text NOT NULL UNIQUE,
  password text NOT NULL
);

CREATE TABLE drawings (
  id SERIAL PRIMARY KEY,
  nanoid text NOT NULL UNIQUE,
  creator_id INTEGER REFERENCES users (id),
  title text NOT NULL
);
