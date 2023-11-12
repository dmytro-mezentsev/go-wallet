CREATE TABLE IF NOT EXISTS wallets (
                                      id uuid PRIMARY KEY,
                                      amount numeric(20, 4) NOT NULL DEFAULT 0
)