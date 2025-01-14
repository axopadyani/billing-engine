CREATE TABLE IF NOT EXISTS loans (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    amount NUMERIC NOT NULL,
    payment_duration_weeks INTEGER NOT NULL,
    payment_amount NUMERIC NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX ON loans(user_id);
