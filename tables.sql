create table "users" (
    id           serial PRIMARY KEY,
    created_at   timestamp default current_timestamp,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    username     text NOT NULL UNIQUE,
    password     text NOT NULL,
    email        text NOT NULL UNIQUE
);

CREATE TABLE "tokens" (
    id          serial PRIMARY KEY,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    key text    NOT NULL UNIQUE,
    user_id     integer,
    expires_at  bigint
);
