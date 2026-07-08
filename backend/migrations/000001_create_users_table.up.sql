CREATE TABLE users
(
    id            uuid                     DEFAULT gen_random_uuid() NOT NULL
        PRIMARY KEY,
    email         varchar(100)                                       NOT NULL
        UNIQUE,
    username      varchar(50)                                        NOT NULL
        UNIQUE,
    password_hash varchar(60)                                        NOT NULL,
    created_at    timestamp with time zone DEFAULT now()             NOT NULL,
    updated_at    timestamp with time zone DEFAULT now()             NOT NULL
);