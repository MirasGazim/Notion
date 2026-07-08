CREATE TABLE blocks
(
    id           uuid                     DEFAULT gen_random_uuid() NOT NULL
        PRIMARY KEY,
    type         varchar                                            NOT NULL,
    parent_id    uuid
        REFERENCES blocks
            ON DELETE CASCADE,
    content      jsonb                    DEFAULT '{}'::jsonb       NOT NULL,
    position     double precision         DEFAULT 0                 NOT NULL,
    workspace_id uuid                                               NOT NULL
        CONSTRAINT fk_blocks_workspace
            REFERENCES workspaces
            ON DELETE CASCADE,
    created_by   uuid                                               NOT NULL
        CONSTRAINT fk_blocks_created_by
            REFERENCES users,
    created_at   timestamp with time zone DEFAULT now()             NOT NULL,
    updated_at   timestamp with time zone DEFAULT now()             NOT NULL
);