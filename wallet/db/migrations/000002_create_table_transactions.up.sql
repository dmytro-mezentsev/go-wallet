CREATE TABLE transactions
(
    id                      UUID PRIMARY KEY NOT NULL,
    user_id                 UUID             NOT NULL,
    wallet_id               UUID             NOT NULL,
    transaction_type        TEXT             NOT NULL,
    amount                  numeric(20, 4)   NOT NULL,
    amount_before           numeric(20, 4)   NOT NULL,
    amount_after            numeric(20, 4)   NOT NULL,
    from_payment_system     TEXT             NOT NULL,
    from_payment_identifier TEXT             NOT NULL,
    to_payment_system       TEXT             NOT NULL,
    to_payment_identifier   TEXT             NOT NULL,
    currency                TEXT             NOT NULL,
    description             TEXT,
    created_at              DATE             NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets (id)
);

CREATE INDEX transactions_wallet_id_idx ON transactions (wallet_id);
CREATE INDEX transactions_user_id_idx ON transactions (user_id);
