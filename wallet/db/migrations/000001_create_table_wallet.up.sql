CREATE TABLE IF NOT EXISTS wallet (
                                      id uuid PRIMARY KEY,
                                      user_id uuid,
                                      balance numeric(20, 2) NOT NULL DEFAULT 0,
)