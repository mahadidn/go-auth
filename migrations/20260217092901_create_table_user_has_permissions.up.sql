CREATE TABLE user_has_permissions (
    user_id       BINARY(16) NOT NULL,
    permission_id BINARY(16) NOT NULL,

    CONSTRAINT pk_user_has_permissions            PRIMARY KEY (user_id, permission_id),
    CONSTRAINT fk_user_has_permissions_user       FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_user_has_permissions_permission FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;