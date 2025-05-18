CREATE TABLE charges (
    fingerprint      TEXT PRIMARY KEY,  -- for idempotency
    local_currency   TEXT,
    local_amount     REAL,
    platform         TEXT,
    instrument       TEXT,
    timestamp        INT,  -- Unix seconds
    charge_currency  TEXT,
    charge_amount    REAL,
    purpose          TEXT,
    ingested_at      INT  -- Unix seconds
);
