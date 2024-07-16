CREATE TABLE users (
    id bigserial,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    timezone varchar(255) NOT NULL DEFAULT 'UTC',
    inbox_id bigint NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE tokens (
    id bigserial,
    user_id bigint NOT NULL,
    token varchar(255) NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE projects (
    id bigserial,
    user_id bigint NOT NULL,
    content text NOT NULL,
    description text NOT NULL DEFAULT '',
    parent_id bigint,
    child_order bigint NOT NULL DEFAULT 0,
    inbox boolean NOT NULL DEFAULT FALSE,
    favorite boolean NOT NULL DEFAULT FALSE,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamp,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES projects (id) ON DELETE CASCADE,
);

CREATE TYPE due AS (
    data timestamp,
    recurring boolean
);

CREATE TYPE duration AS (
    amount int,
    unit varchar ( 255));

CREATE TABLE tasks (
    id bigserial,
    user_id bigint NOT NULL,
    content text NOT NULL,
    description text NOT NULL DEFAULT '',
    due due NOT NULL DEFAULT ROW (NULL, NULL),
    duration duration NOT NULL DEFAULT ROW (NULL, NULL),
    priority int NOT NULL DEFAULT 0,
    project_id bigint,
    parent_id bigint,
    child_order bigint NOT NULL DEFAULT 0,
    done boolean NOT NULL DEFAULT FALSE,
    done_at timestamp,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamp,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES tasks (id) ON DELETE CASCADE,
);

CREATE TABLE labels (
    id bigserial,
    user_id bigint NOT NULL,
    name varchar(255) NOT NULL,
    color varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
);

