CREATE TABLE calculations (
                              id         UUID PRIMARY KEY,
                              num_a      FLOAT8      NOT NULL,
                              num_b      FLOAT8      NOT NULL,
                              sign       VARCHAR(1)  NOT NULL,
                              result     FLOAT8      NOT NULL,
                              created_at TIMESTAMPTZ NOT NULL
);