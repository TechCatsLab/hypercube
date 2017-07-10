CREATE DATABASE IF NOT EXISTS core;
SET DATABASE = core;

CREATE TABLE IF NOT EXISTS message (
  messageid     SERIAL        PRIMARY KEY,
  source          INT,
  target          INT         NOT NULL,
  type            INT         NOT NULL,
  issend          BOOL        DEFAULT FALSE,
  content         SRTING(256),
  created         TIMESTAMP   DEFAULT current_timestamp()
);
