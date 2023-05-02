CREATE TABLE user_account(
    id serial not null unique,
    user_name varchar (50) not null unique,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);
CREATE TABLE habit (
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255)
);
CREATE TABLE user_habit (
    id serial not null unique,
    user_id int references user_account (id) on delete cascade not null,
    habit_id int references habit (id) on delete cascade not null
);
CREATE TABLE habit_tracker (
    id serial not null unique,
    user_habit_id int references user_habit (id) on delete cascade not null,
    unit_of_messure varchar(50) not null,
    goal varchar(50) not null,
    frequency varchar(255) not null,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE NOT NULL,
    counter NUMERIC(10, 2),
    done boolean not null default false
);

CREATE TABLE reward (
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE user_reward {
    id serial not null unique,
    user_id int references user_account (id) on delete cascade not null,
    reward_id int references reward (id) on delete cascade not null,
    habit_id int references habit (id) on delete cascade not null
}