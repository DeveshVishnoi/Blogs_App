alter table users_invitation
add column expiry timestamp(0) with time zone not null;