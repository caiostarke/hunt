CREATE TABLE users (
    id UUID PRIMARY KEY,
    level INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    profile_picture TEXT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password CHAR(60) NOT NULL
);
