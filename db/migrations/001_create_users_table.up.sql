CREATE TABLE IF NOT EXISTS users  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL, -- for auth, if needed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
