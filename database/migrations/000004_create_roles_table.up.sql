CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    role_id INTEGER REFERENCES roles(id)
);
