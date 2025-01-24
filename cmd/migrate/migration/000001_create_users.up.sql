create extension if not exists citext;

Create table if not exists users(
    id bigserial primary key,
    email citext unique not null,
    username varchar(255) unique not null,
    password bytea not null,
    created_at timestamp(0) with time zone not null default Now()

)