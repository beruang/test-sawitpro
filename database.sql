/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users
(
    id        serial PRIMARY KEY,
    slug      char(20)           not null,
    full_name varchar(60)        not null,
    phone     varchar(15) unique not null,
    password  char(60)           not null,
)