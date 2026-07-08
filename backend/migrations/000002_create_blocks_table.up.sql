CREATE TABLE workspaces
(
    id         uuid                     DEFAULT gen_random_uuid() NOT NULL
        PRIMARY KEY,
    owner_id   uuid                                               NOT NULL
        REFERENCES users
            ON DELETE CASCADE,
    name       varchar                                            NOT NULL,
    created_at timestamp with time zone DEFAULT now()             NOT NULL
);