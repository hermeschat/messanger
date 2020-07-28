CREATE TABLE IF NOT EXISTS users
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    phone_number VARCHAR(200) NOT NULL,
    created_at   timestamp DEFAULT NOW(),
    updated_at   timestamp DEFAULT NOW(),
    deleted_at   timestamp DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS roles
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(200) NOT NULL,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS user_role
(
    user_id    INT,
    role_id    INT,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (role_id) REFERENCES roles (id),
    PRIMARY KEY (user_id, role_id),
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
)