create table if not exists Services(
    id serial not null ,
    service varchar(60) not null ,
    login   varchar(60) not null ,
    password varchar(60) not null ,
    user_id integer not null ,
    primary key (id)
    )