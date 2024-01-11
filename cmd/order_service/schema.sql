-- create database if not exists bond_trading ENCODING = 'UTF8';

drop type if exists orderaction cascade;
create type orderaction as enum ('BUY', 'SELL', 'CANCEL');
alter type orderaction owner to postgres;

drop type if exists orderstatus cascade;
create type orderstatus as enum ('PENDING', 'OPEN', 'FILLED', 'CANCELLED');
alter type orderstatus owner to postgres;

create table if not exists users
(
    id            bigint                   not null constraint users_pk primary key,
    username      varchar                  not null,
    email         varchar                  not null,
    password_hash varchar                  not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone,
    active        boolean                  not null
);

alter table users owner to postgres;

create table if not exists orders
(
    id          varchar                  not null constraint orders_pk primary key,
    user_id     bigint                   not null constraint orders_users_fk references users,
    bond_id     varchar                  not null,
    quantity    integer                  not null,
    filled      money                    not null,
    price       money                    not null,
    action      orderaction              not null,
    status      orderstatus              not null,
    expiress_at timestamp with time zone not null,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone
);

alter table orders owner to postgres;
