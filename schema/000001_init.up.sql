CREATE TABLE user (
    user_id serial not null unique,
    user_name varchar(50) not null unique,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    e - mail varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE habbit (
    habbit_id serial not null unique,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE user_habbit (
    user_habbit_id serial not null unique,
    user_id int references user (id) on delete cascade not null,
    habbit_id int references habbit (id) on delete cascade not null
);

CREATE TABLE habbit_tracker (
    habbit_tracker_id serial not null unique,
    user_habbit_id serial not null,
    unit_of_messure varchar(50) not null,
    goal varchar(50) not null,
    frequency varchar(255) not null,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE NOT NULL,
    counter NUMERIC(10,2)
);

CREATE TABLE reward (
    reward_id serial not null unique,
    habbit_tracker_id int references habbit_tracker (id) on delete cascade not null,
    title varchar(255) not null,
    description varchar(255)
);