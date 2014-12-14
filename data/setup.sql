
drop table users;

create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null   
);

drop table sessions;

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null   
);

drop table conversations;

create table conversations (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  topic      text,
  user_id    integer references users(id),
  created_at timestamp not null       
);
drop table replies;

create table replies (
  id              serial primary key,
  uuid            varchar(64) not null unique,
  body            text,
  user_id         integer references users(id),
  conversation_id integer references conversations(id),
  created_at      timestamp not null  
);