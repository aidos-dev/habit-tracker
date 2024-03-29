CREATE TABLE user_account(
    id serial not null unique,
    user_name varchar (50) unique,
    tg_user_name varchar (50) unique,
    first_name varchar(50),
    last_name varchar(50),
    email varchar(255) unique,
    password_hash varchar(255),
    role varchar (50) DEFAULT 'user_basic'
);
CREATE TABLE habit (
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE habit_tracker (
    id serial not null unique,
    habit_id NUMERIC(10),
    unit_of_messure varchar(50),
    goal varchar(50),
    frequency varchar(255),
    start_date DATE DEFAULT CURRENT_DATE,
    end_date DATE,
    counter NUMERIC(10, 2),
    done boolean DEFAULT false
);

CREATE TABLE user_habit (
    id serial not null unique,
    user_id int references user_account (id) ON DELETE CASCADE not null,
    habit_id int references habit (id) ON DELETE CASCADE not null,
    habit_tracker_id int,
    UNIQUE (user_id, habit_id, habit_tracker_id),
    PRIMARY KEY (user_id, habit_id)
);


CREATE TABLE reward (
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255),
    UNIQUE (title, description)
);

CREATE TABLE user_reward (
    id serial not null unique,
    user_id int references user_account (id) ON DELETE CASCADE not null,
    habit_id int references habit (id) ON DELETE CASCADE not null,
    reward_id int references reward (id) ON DELETE CASCADE not null,
    UNIQUE (user_id, habit_id, reward_id)
);





