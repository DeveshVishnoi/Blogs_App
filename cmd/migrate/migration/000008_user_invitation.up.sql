create table if not EXISTS users_invitation(
    token bytea primary key,
    user_id bigint not null
);