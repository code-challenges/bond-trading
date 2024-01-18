-- create database if not exists bond_trading ENCODING = 'UTF8';

drop type if exists orderaction cascade;
create type orderaction as enum ('BUY', 'SELL', 'CANCEL');
alter type orderaction owner to postgres;

drop type if exists orderstatus cascade;
create type orderstatus as enum ('PENDING', 'OPEN', 'FILLED', 'CANCELLED');
alter type orderstatus owner to postgres;

create table if not exists users
(
    id            serial            not null constraint users_pk primary key,
    email         varchar           not null,
    password_hash varchar           not null,
    active        boolean           not null,
    created_at    timestamptz       not null,
    updated_at    timestamptz
);
alter table users owner to postgres;

create table if not exists orders
(
    id          varchar             not null constraint orders_pk primary key,
    user_id     bigint              not null constraint orders_users_fk references users,
    bond_id     integer             not null,
    quantity    integer             not null,
    filled      integer             not null,
    price       decimal(9, 4)       not null,
    action      orderaction         not null,
    status      orderstatus         not null,
    expires_at  timestamptz         not null,
    created_at  timestamptz         not null,
    updated_at  timestamptz
);
alter table orders owner to postgres;

-- create dev user with password: 123
-- insert into users (id, email, password_hash, created_at, active)
-- values (1, 'dev@localhost', '$2a$10$bmPeLZ/DPG.7f7PbcoC.g.F0jZ4KPl64Tr4p8kxNJ/jRz0LFwPe9K', now(), true);
