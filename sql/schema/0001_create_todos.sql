-- +goose Up

create table todos (
    id UUID primary key,
    description text not null,
    done boolean not null
);

-- +goose Down

drop table todos;