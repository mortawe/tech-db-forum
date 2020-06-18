CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    nickname CITEXT PRIMARY KEY NOT NULL,
    email    CITEXT UNIQUE      NOT NULL,
    about    TEXT               NOT NULL,
    fullname TEXT               NOT NULL
);

CREATE UNIQUE INDEX ON users (nickname, email);
CREATE UNIQUE INDEX ON users (nickname, email, about, fullname);
-- CREATE UNIQUE INDEX ON users (nickname ASC); ALREADY EXISTS
CREATE UNIQUE INDEX ON users (nickname DESC);

CREATE UNLOGGED TABLE forums
(
    slug     CITEXT PRIMARY KEY                                   NOT NULL,
    title    TEXT                                                 NOT NULL,
    nickname CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    posts    INTEGER DEFAULT 0                                    NOT NULL,
    threads  INTEGER DEFAULT 0                                    NOT NULL
);

-- CREATE UNIQUE INDEX ON forums(slug); ALREADY EXISTS ?
-- CREATE UNIQUE INDEX ON forums(slug, title, nickname, posts, threads); DOES NOT WORKS BECAUSE OF DEADLOCK =)

CREATE UNLOGGED TABLE forum_users
(
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug   CITEXT REFERENCES forums (slug) ON DELETE CASCADE    NOT NULL,
    PRIMARY KEY (slug, author)
);

CREATE INDEX ON forum_users (slug);
CREATE INDEX ON forum_users (author);

CREATE UNLOGGED TABLE threads
(
    author     CITEXT REFERENCES users (nickname) ON DELETE CASCADE  NOT NULL,
    created    TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_slug CITEXT REFERENCES forums (slug) ON DELETE CASCADE     NOT NULL,
    id         SERIAL PRIMARY KEY                                    NOT NULL,
    message    TEXT                                                  NOT NULL,
    slug       CITEXT,
    title      TEXT                                                  NOT NULL,
    votes      INTEGER                     DEFAULT 0                 NOT NULL
);

CREATE INDEX ON threads(slug, author);
CREATE INDEX ON threads(forum_slug, created ASC);
CREATE INDEX ON threads(forum_slug, created DESC);
-- CREATE INDEX ON threads(id); ALREADY EXISTS
CREATE INDEX ON threads(slug, id);
CREATE INDEX ON threads(id, forum_slug);
CREATE INDEX ON threads(slug, id, forum_slug);

CREATE UNLOGGED TABLE posts
(
    author     CITEXT REFERENCES users (nickname) ON DELETE CASCADE  NOT NULL,
    created    TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_slug CITEXT REFERENCES forums (slug) ON DELETE CASCADE     NOT NULL,
    id         SERIAL PRIMARY KEY                                    NOT NULL,
    edited     BOOL                        DEFAULT 'false'           NOT NULL,
    message    TEXT                                                  NOT NULL,
    parent     INTEGER                                               NOT NULL,
    thread     INTEGER REFERENCES threads (id) ON DELETE CASCADE     NOT NULL,
    path       INTEGER ARRAY               DEFAULT '{}'              NOT NULL
);

CREATE UNIQUE INDEX ON posts(id, thread);
CREATE UNIQUE INDEX ON posts(id, author);
CREATE INDEX ON posts(thread, path DESC);
CREATE INDEX ON posts(thread, path ASC);
CREATE INDEX ON posts(thread, id DESC);
CREATE INDEX ON posts(thread, id ASC);

CREATE UNLOGGED TABLE votes
(
    nickname  CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    thread_id INTEGER REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    vote      SMALLINT                                             NOT NULL,
    PRIMARY KEY (thread_id, nickname)
);

CREATE UNIQUE INDEX ON votes(thread_id, nickname);

CREATE INDEX ON threads (slug, id);

-- PATH TO POST UPDATE
CREATE FUNCTION update_path() RETURNS TRIGGER AS
$$
DECLARE
    temp INT ARRAY;
    t    INTEGER;
BEGIN
    IF new.parent ISNULL OR new.parent = 0 THEN
        new.path = ARRAY [new.id];
    ELSE
        SELECT thread
        INTO t
        FROM posts
        WHERE id = new.parent;
        IF t ISNULL OR t <> new.thread THEN
            RAISE EXCEPTION 'Not in this thread ID ' USING HINT = 'Please check your parent ID';
        END IF;

        SELECT path
        INTO temp
        FROM posts
        WHERE id = new.parent;
        new.path = array_append(temp, new.id);

    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

-- VOTE VALUE UPDATE
CREATE FUNCTION vote_count_upd() RETURNS TRIGGER AS
$$
BEGIN
    IF (old.vote != new.vote) THEN
        UPDATE threads
        SET votes = (votes - old.vote + new.vote)
        WHERE id = new.thread_id;
    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_upd
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_upd();

CREATE FUNCTION vote_count_insert() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads
    SET votes = (votes + new.vote)
    WHERE id = new.thread_id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_insert
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_insert();


-- UPDATE FORUM_USER TABLE AFTER INSERTS
CREATE FUNCTION insert_forum_user_from_threads_or_psoots() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO forum_users
    VALUES (new.author, new.forum_slug)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_user_from_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();

CREATE TRIGGER update_forum_user_from_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();

-- UPDATE POSTS AND THREADS COUNTERS IN FORUMS
CREATE FUNCTION update_forum_counter_posts() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_post_insert
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_posts();

CREATE FUNCTION update_forum_counter_threads() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_thread_insert
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_threads();

-- НАБОР 1
-- 05:52:52.176 INFO Requests per second: 1831.10
-- 05:53:02.177 INFO Requests per second: 1914.50
-- 05:53:12.177 INFO Requests per second: 1795.80
-- 05:53:22.177 INFO Requests per second: 1705.30
-- 05:53:32.177 INFO Requests per second: 1769.50
-- 05:53:42.178 INFO Requests per second: 1720.20
-- 05:53:52.178 INFO Requests per second: 1531.10
-- 05:54:02.179 INFO Requests per second: 1131.10


-- НАБОР 2 (текущий)
-- 14:39:08.907 INFO Requests per second: 2736.00
-- 14:39:18.907 INFO Requests per second: 2990.90
-- 14:39:28.907 INFO Requests per second: 2756.10
-- 14:39:38.907 INFO Requests per second: 2451.80
-- 14:39:48.907 INFO Requests per second: 2643.80
-- 14:39:58.907 INFO Requests per second: 2788.70
-- 14:40:08.907 INFO Requests per second: 2683.30
-- 14:40:18.908 INFO Requests per second: 2618.00
-- 14:40:28.908 INFO Requests per second: 2632.60

VACUUM ANALYSE;

-- НАБОР 3 (с ANALYSE)
-- 14:43:43.983 INFO Requests per second: 3047.30
-- 14:43:53.983 INFO Requests per second: 3021.50
-- 14:44:03.983 INFO Requests per second: 3069.80
-- 14:44:13.983 INFO Requests per second: 2952.90