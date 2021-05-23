create table if not exists users
(
    id            serial       not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    first_name    varchar(255) not null,
    last_name     varchar(255) not null,
    registered_at timestamp    not null default now()
);

create table if not exists refreshSessions
(
    id           serial                    not null unique,
    userId       int references users (id) not null,
    refreshToken varchar(255)              not null,
    expiresIn    timestamp                 not null default now(),
    createdAt    timestamp                 not null default now()
);

create table if not exists admins
(
    id            serial       not null unique,
    login         varchar(255) not null,
    password_hash varchar(255) not null
);

create table if not exists adminsRefreshSessions
(
    id           serial                     not null unique,
    adminId      int references admins (id) not null,
    refreshToken varchar(255)               not null,
    expiresIn    timestamp                  not null default now(),
    createdAt    timestamp                  not null default now()
);

create table if not exists categories
(
    id              serial       not null unique,
    category        varchar(255) not null,
    parent_category int references categories (id)
);

create table if not exists contacts_info
(
    id           serial       not null unique,
    name         varchar      not null,
    phone_number varchar(255) not null,
    email        varchar(255) not null,
    location     varchar(255) not null default 'Kiev'
);

create table if not exists ads
(
    id          serial                                              not null unique,
    userId      int references users (id) on delete cascade         not null,
    title       varchar(255)                                        not null,
    category_id int references categories (id)                      not null,
    description text                                                not null,
    price       int                                                 not null,
    contacts_id int references contacts_info (id) on delete cascade not null,
    published   boolean                                             not null default false,
    images_url  varchar[]                                           not null
);

insert into categories(category)
values ('Транспорт'); -- 1
insert into categories(category, parent_category)
values ('Автобусы', 1),
       ('Водный транспорт', 1),
       ('Мототехника', 1),
       ('Прицепы', 1);
insert into categories(category)
values ('Одежда и обувь'); -- 6
insert into categories(category, parent_category)
values ('Аксессуары', 6),
       ('Женская обувь', 6),
       ('Женская одежда', 6),
       ('Мужская обувь', 6),
       ('Другое', 6);
insert into categories(category)
values ('Недвижимость'); -- 12
insert into categories(category, parent_category)
values ('Аренда домов', 12),
       ('Аренда квартир', 12),
       ('Продажа квартир', 12),
       ('Продажа домов', 12);

insert into admins (login, password_hash)
values ('admin1@gmail.com',
        '6473667364666d31336b326d444d4b53464d53444b469347bb6d872adfd5180cbab2fa8fd64014064e45'); --password - dfghjk1503
insert into admins (login, password_hash)
values ('admin2@gmail.com',
        '6473667364666d31336b326d444d4b53464d53444b469f8d826c82de9b38b0572acbd7cc465bcc424831'); --password - dfghjk

---- help function
create or replace function make_tsvector(title varchar, description text)
    returns tsvector as $$
begin
    RETURN (setweight(to_tsvector('russian', title), 'A') || setweight(to_tsvector('russian', description), 'B'));
end
    $$ language 'plpgsql' immutable;

---- gin index on ads table
create index if not exists idx_fts_ads on ads
    using gin(make_tsvector(title, description));