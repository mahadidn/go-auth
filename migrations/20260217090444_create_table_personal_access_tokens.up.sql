CREATE TABLE personal_access_tokens (
    id            BINARY(16)   NOT NULL, 
    token_hash    VARCHAR(64)  NOT NULL, 
    user_id       BINARY(16)   NOT NULL,
    token_name    VARCHAR(100) NOT NULL,
    last_used_at  TIMESTAMP    NULL,       
    expires_at    TIMESTAMP    NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_personal_access_tokens      PRIMARY KEY (id),
    CONSTRAINT fk_personal_access_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT idx_token_hash                 UNIQUE (token_hash) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;