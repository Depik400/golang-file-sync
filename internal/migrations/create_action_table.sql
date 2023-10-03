drop table if exists system_actions;
create table if not exists system_actions (
    id integer primary key not null,
    directory_name varchar(255) not null,
    file_name varchar(255) not null,
    action_key varchar(255) not null,
    created_at integer default(unixepoch()),
    unique(directory_name, file_name, action_key)
)