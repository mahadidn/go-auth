CREATE TABLE user_has_roles (
    user_id BINARY(16) NOT NULL,
    role_id BINARY(16) NOT NULL,

    CONSTRAINT pk_user_has_roles      PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_has_roles_user FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_user_has_roles_role FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;