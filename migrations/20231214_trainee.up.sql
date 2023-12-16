create table users (
    id serial primary key,
    username varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp not null default now()
);

create table slugs (
    id serial primary key,
    name varchar(255) not null unique,
    created_at timestamp not null default now()
);

create table membership (
    user_id int not null,
    slug_id int not null,
    created_at timestamp not null default now(),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (slug_id) references slugs(id) on delete cascade,
    primary key (user_id, slug_id)
);