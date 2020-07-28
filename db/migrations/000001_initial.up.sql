CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL PRIMARY KEY,
    phone_number VARCHAR(200)  NOT NULL,
    username     varchar(200)  NOT NULL,
    password     varchar(1000) NOT NULL,
    created_at   timestamp DEFAULT NOW(),
    updated_at   timestamp DEFAULT NOW(),
    deleted_at   timestamp DEFAULT NULL
-- unique index on username
);

CREATE TABLE IF NOT EXISTS roles
(
    id         SERIAL PRIMARY KEY,
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
);

--- channels
CREATE TABLE IF NOT EXISTS channels
(
    id         SERIAL PRIMARY KEY,
    creator_id INT,
    FOREIGN KEY (creator_id) references users (id),
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS channel_members
(
    channel_id INT,
    user_id    INT,
    FOREIGN KEY (channel_id) REFERENCES channels (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (channel_id, user_id),
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS channel_permissions
(
    channel_id INT,
    user_id    INT,
    FOREIGN KEY (channel_id) REFERENCES channels (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (channel_id, user_id),
    permission INT,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
);
--- messages table
CREATE TABLE IF NOT EXISTS messages
(
    id         SERIAL PRIMARY KEY,
    origin_id  INT,
    dst_id     INT,
    parent_id  INT,
    body       TEXT,
    state      INT,
    FOREIGN KEY (origin_id) REFERENCES users (id),
    FOREIGN KEY (dst_id) REFERENCES channels (id),
    FOREIGN KEY (parent_id) REFERENCES messages (id),
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp DEFAULT NULL
);