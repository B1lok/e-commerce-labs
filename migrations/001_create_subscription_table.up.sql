CREATE TYPE frequency_enum AS ENUM ('daily', 'hourly');

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    frequency frequency_enum NOT NULL,
    token VARCHAR(36) NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (email, city, frequency)
);