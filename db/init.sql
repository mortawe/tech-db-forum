CREATE EXTENSION IF NOT EXISTS citext;

CREATE unlogged TABLE users
(
    about    varchar not null,
    email    citext unique not null,
    fullname citext not null,
    nickname citext unique not null
);


CREATE unlogged TABLE forums
(
    slug     citext primary key,
    title    varchar not null,
    nickname citext references users (nickname) on delete cascade,
    posts    integer not null default 0,
    threads  integer not null default 0
);

CREATE unlogged TABLE threads
(
    author  citext not null,
    created timestamp with time zone,
    forum   citext not null references forums (slug) on delete cascade,
    id      serial primary key,
    message varchar not null,
    slug    citext,
    title   varchar not null,
    votes   integer not null default 0
);

CREATE unlogged TABLE posts
(
    author  citext not null references users (nickname) on delete cascade,
    created timestamp with time zone,
    forum   citext not null references forums (slug) on delete cascade ,
    id      serial,
    edited  bool default 'false',
    message varchar not null,
    parent  integer not null,
    thread  integer not null,
    path integer array default '{}'
);

CREATE unlogged TABLE votes
(
    nickname citext not null references users (nickname) on delete cascade ,
    threadID integer not null references threads (id) on delete cascade ,
    vote     integer not null
);



CREATE FUNCTION update_path()
RETURNS trigger AS
$BODY$
DECLARE
    temp int array;
BEGIN
    IF NEW.parent = 0 THEN
        NEW.path = ARRAY[NEW.id];
    ELSE
        SELECT path
        INTO temp
        FROM posts
        WHERE
            id = NEW.parent;
        NEW.path = array_append(temp, NEW.id);
    END IF;
    return new;
end;
$BODY$ LANGUAGE plpgsql;



create trigger update_posts_path
    before insert
    on posts
    for each row
execute procedure update_path();



CREATE unlogged TABLE user_forum
(
    author  citext not null references users( nickname) on DELETE CASCADE,
    forum citext not null  references forums(slug) on delete cascade
);

create function update_threads_count() returns trigger
    language plpgsql
as
$$
BEGIN
    update forums
    set threads = threads + 1
    where NEW.forum = slug;

    insert into user_forum
    values (new.author, new.forum);

    return new;
end;
$$;

alter function update_threads_count() owner to postgres;

create function update_posts_count() returns trigger
    language plpgsql
as
$$
BEGIN
    update forums
    set posts = posts + 1
    where NEW.forum = slug;

    insert into user_forum
    values (new.author, new.forum);

    return new;
end;
$$;

alter function update_posts_count() owner to docker;

create trigger update_threads_count
    after insert
    on threads
    for each row
execute procedure update_threads_count();


create trigger update_posts_count
    after insert
    on posts
    for each row
execute procedure update_posts_count();

create index on posts (thread);
create index ON posts using gin (path);
create index on threads (slug);
create index on threads (created, forum);
create index on user_forum (forum);
