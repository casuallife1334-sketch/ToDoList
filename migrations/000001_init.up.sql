CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
	description VARCHAR(300),
	completed BOOLEAN,
	created_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP
);