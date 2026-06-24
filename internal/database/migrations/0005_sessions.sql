-- +goose up
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

-- +goose down
DROP TABLE sessions;
DROP INDEX sessions_expiry_idx;
