CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id),
    bio TEXT DEFAULT '',
    avatar_url TEXT DEFAULT ''
);
