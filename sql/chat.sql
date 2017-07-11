CREATE DATABASE IF NOT EXISTS core;
SET DATABASE = core;

CREATE TABLE IF NOT EXISTS message (
  messageid       SERIAL            PRIMARY KEY,
  source          STRING(128),
  target          STRING(128)     NOT NULL,
  type            INT             NOT NULL,
  issend          BOOL            DEFAULT FALSE,
  content         STRING(256),
  created         TIMESTAMP       DEFAULT current_timestamp()
);
