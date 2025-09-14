CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(50),
  email VARCHAR(50) NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX users_email_unique_idx ON users(email);