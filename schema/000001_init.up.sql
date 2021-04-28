create table users
(
    id            serial       not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    first_name    varchar(255) not null,
    last_name     varchar(255) not null,
    registered_at timestamp    not null default now()
);

create table refreshSessions
(
    id           serial                    not null unique,
    userId       int references users (id) not null,
    refreshToken varchar(255)              not null,
    ua           varchar(255),
    ip           varchar(255),
    expiresIn    timestamp                 not null default now(),
    createdAt    timestamp                 not null default now()
);

create table admin_users
(
    id            serial       not null unique,
    login         varchar(255) not null,
    password_hash varchar(255) not null
);

create table categories
(
    id              serial       not null unique,
    category        varchar(255) not null,
    parent_category int references categories (id)
);

create table contacts_info
(
    id           serial       not null unique,
    name         varchar      not null,
    phone_number varchar(255) not null,
    email        varchar(255) not null,
    location     varchar(255) not null default 'Kiev'
);

create table ads
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

