create table if not exists customer
    (id serial primary key,
    email varchar(100) not null unique
    );