BEGIN;

CREATE TABLE IF NOT EXISTS outbox
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    topic      TEXT                                  NOT NULL,
    key        BYTEA                                 NOT NULL,
    value      BYTEA                                 NOT NULL
);

COMMIT;