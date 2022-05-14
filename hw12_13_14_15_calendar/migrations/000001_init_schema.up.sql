create table events (
    id                 bigserial primary key,
    title              text not null,
    start_event        timestamp not null,
    end_event          timestamp not null,
    description        text,
    id_user            bigint not null,
    notification       timestamp,
    notificationStatus int default 0
);

create table users (
    id              bigserial primary key,
    nickname        text not null default 'test_user'
);

alter table events add foreign key (id_user) references users (id);

create index on events (id_user);

insert into users (nickname) values ('Alice'), ('Bob');

insert into events (title, start_event, end_event, description, id_user, notification, notificationStatus)
 values ('default_title', '2022-07-01 06:30:30', '2022-08-01 06:30:30', 'test_event', 1, '2022-07-01 06:30:30', 0);
