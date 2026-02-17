CREATE TABLE users (
    id            BINARY(16)       NOT NULL,
    username      VARCHAR(45)      NOT NULL,
    email         VARCHAR(255)     NOT NULL,
    password      VARCHAR(255)     NOT NULL,
    created_at    TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP
                                            ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk_users           PRIMARY KEY (id),
    CONSTRAINT uq_users_username  UNIQUE      (username),
    CONSTRAINT uq_users_email     UNIQUE      (email)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;