CREATE TABLE roles(
    id SERIAL PRIMARY KEY,
    role_type INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- CREATE TABLE accesses (
--     id SERIAL PRIMARY KEY,
--     role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
--     access_type INTEGER NOT NULL,
--     created_at TIMESTAMPTZ DEFAULT now()
-- );

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    title TEXT,
    statement TEXT,
    owner INTEGER REFERENCES users(id) ON DELETE SET NULL,
    input TEXT,
    output TEXT,
    memory_limit INTEGER,
    time_limit INTEGER,
    state INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    code TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    state INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

