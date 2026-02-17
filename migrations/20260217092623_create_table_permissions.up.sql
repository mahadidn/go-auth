CREATE TABLE permissions (
    id          BINARY(16)   NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
                                      ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT pk_permissions      PRIMARY KEY (id),
    CONSTRAINT uq_permissions_name UNIQUE      (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;