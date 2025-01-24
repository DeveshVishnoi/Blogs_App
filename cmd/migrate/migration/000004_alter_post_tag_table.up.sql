alter table posts add column tags varchar(100) [];

alter table posts add column updated_At timestamp(0) with time zone NOT NULL DEFAULT NOW()