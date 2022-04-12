create table events (
    id                 bigserial primary key,
    title              text not null,
    start_event        timestamp not null,
    end_event          timestamp not null,
    description        text,
    id_user            bigint not null,
    notification       timestamp,
    notificationSended boolean default false
);

create table users (
    id              bigserial primary key,
    nickname        text not null default 'test_user'
);

alter table events add foreign key (id_user) references users (id);

create index on events (id_user);

insert into users (nickname) values ('Alice'), ('Bob');
