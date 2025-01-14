CREATE TABLE IF NOT EXISTS loan_payments (
    id UUID PRIMARY KEY,
    loan_id UUID NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id)
);
