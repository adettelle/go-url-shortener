create table if not exists url_mapping
    (id serial primary key,
    short_url varchar(30),
    original_url varchar(500)
    );