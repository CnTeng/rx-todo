CREATE TABLE users (
  id bigserial,
  username varchar(255) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  email varchar(255) NOT NULL UNIQUE,
  timezone varchar(255) NOT NULL DEFAULT 'UTC',
  inbox_id bigint NOT NULL DEFAULT 0,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE TABLE tokens (
  id bigserial,
  user_id bigint NOT NULL,
  token varchar(255) NOT NULL UNIQUE,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE projects (
  id bigserial,
  user_id bigint NOT NULL,
  name text NOT NULL,
  description text NOT NULL DEFAULT '',
  child_order bigint NOT NULL UNIQUE,
  inbox boolean NOT NULL DEFAULT FALSE,
  favorite boolean NOT NULL DEFAULT FALSE,
  archived boolean NOT NULL DEFAULT FALSE,
  archived_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  UNIQUE (user_id, child_order)
);

CREATE TYPE due AS (date timestamp, recurring boolean);

CREATE TYPE duration AS (amount int, unit varchar(255));

CREATE TABLE tasks (
  id bigserial,
  user_id bigint NOT NULL,
  name text NOT NULL,
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
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
  FOREIGN KEY (parent_id) REFERENCES tasks (id) ON DELETE CASCADE,
  UNIQUE (project_id, child_order),
  UNIQUE (parent_id, child_order)
);

CREATE TABLE labels (
  id bigserial,
  user_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  color varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  UNIQUE (user_id, name)
);

CREATE TABLE task_labels (
  task_id bigint NOT NULL,
  label_id bigint NOT NULL,
  PRIMARY KEY (task_id, label_id),
  FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE,
  FOREIGN KEY (label_id) REFERENCES labels (id) ON DELETE CASCADE
);

CREATE TABLE reminders (
  id bigserial,
  user_id bigint NOT NULL,
  task_id bigint NOT NULL,
  due due NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
);

CREATE TYPE object_type AS enum('label', 'project', 'reminder', 'task', 'user');

CREATE TABLE deletion_log (
  id bigserial,
  user_id bigint NOT NULL,
  object_type object_type NOT NULL,
  object_id bigint NOT NULL,
  deleted_at timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX deletion_log_deleted_at_idx ON deletion_log (deleted_at);

CREATE VIEW tasks_with_sub_tasks AS
WITH
  task_counts AS (
    SELECT
      parent_id,
      user_id,
      count(*) AS total_tasks,
      count(
        CASE
          WHEN done THEN 1
          ELSE NULL
        END
      ) AS done_tasks
    FROM
      tasks
    GROUP BY
      parent_id,
      user_id
  )
SELECT
  tasks.id,
  tasks.user_id,
  tasks.name,
  tasks.description,
  tasks.due,
  tasks.duration,
  tasks.priority,
  tasks.project_id,
  tasks.parent_id,
  tasks.child_order,
  coalesce(tc.total_tasks, 0) AS total_tasks,
  coalesce(tc.done_tasks, 0) AS done_tasks,
  tasks.done,
  tasks.done_at,
  tasks.archived,
  tasks.archived_at,
  tasks.created_at,
  tasks.updated_at
FROM
  tasks
  LEFT JOIN task_counts tc ON tasks.id = tc.parent_id
  AND tasks.user_id = tc.user_id;

CREATE VIEW projects_with_sub_tasks AS
WITH
  task_counts AS (
    SELECT
      project_id,
      user_id,
      count(*) AS total_tasks,
      count(
        CASE
          WHEN done THEN 1
          ELSE NULL
        END
      ) AS done_tasks
    FROM
      tasks
    GROUP BY
      project_id,
      user_id
  )
SELECT
  projects.id,
  projects.user_id,
  projects.name,
  projects.description,
  projects.child_order,
  projects.inbox,
  projects.favorite,
  coalesce(tc.total_tasks, 0) AS total_tasks,
  coalesce(tc.done_tasks, 0) AS done_tasks,
  projects.archived,
  projects.archived_at,
  projects.created_at,
  projects.updated_at
FROM
  projects
  LEFT JOIN task_counts tc ON projects.id = tc.project_id
  AND projects.user_id = tc.user_id;
