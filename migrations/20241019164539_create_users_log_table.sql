-- +goose Up
create table users_log
(
    id      serial primary key,
    user_id int not null,
    action_time timestamp not null default now(),
    action varchar(50) not null
);

-- +goose Down
drop table users_log;
