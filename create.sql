drop table if exists invite;

create table invite (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT,
    last_name TEXT,
    email TEXT
);

drop table if exists person;

create table person (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    invite INTEGER,
    name TEXT,
    rsvp INTEGER,
    dietary TEXT,
    notes TEXT
);

-- insert into invite (first_name, last_name) values
-- ('dan', 'mitchell'),
-- ('madi', 'mitchell'),
-- ('joe', 'mitchell'),
-- ('rex', 'beck'),
-- ('dan', 'greenfield'),
-- ('tim', 'plouffe');

-- insert into person (invite, name) values
-- (1, 'dan mitchell'),
-- (1, 'alyssa melendez'),
-- (2, 'madi mitchell'),
-- (3, 'joe mitchell');