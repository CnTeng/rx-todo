CREATE TABLE users (
    id bigserial,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    timezone varchar(255) DEFAULT 'UTC',
    PRIMARY KEY (id)
);

CREATE TABLE token (
    id bigserial,
    user_id int NOT NULL,
    token varchar(255) NOT NULL UNIQUE,
    last_used_at timestamp NOT NULL DEFAULT now(),
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE projects (
    id bigserial,
    user_id int NOT NULL,
    project_id int,
    content text NOT NULL,
    description text,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    archive boolean NOT NULL DEFAULT FALSE,
    archive_at timestamp,
    child_order int NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TYPE due AS (
    data timestamp,
    recurring boolean
);

CREATE TABLE tasks (
    id bigserial,
    user_id int NOT NULL,
    content text NOT NULL,
    description text,
    due due,
    duration interval,
    priority int NOT NULL DEFAULT 0,
    project_id int,
    parent_id int,
    child_order int NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    done boolean NOT NULL DEFAULT FALSE,
    done_at timestamp,
    archive boolean NOT NULL DEFAULT FALSE,
    archive_at timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES tasks (id) ON DELETE CASCADE
);

CREATE TABLE labels (
    id bigserial,
    user_id int NOT NULL,
    name varchar(255) NOT NULL,
    color varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
);

