CREATE TABLE role_has_permissions (
    role_id       BINARY(16) NOT NULL,
    permission_id BINARY(16) NOT NULL,

    CONSTRAINT pk_role_has_permissions            PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_has_permissions_role       FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_role_has_permissions_permission FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;