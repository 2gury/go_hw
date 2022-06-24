DROP TRIGGER IF EXISTS score_upd on votes;

DROP TABLE IF EXISTS
    users
CASCADE;

CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username varchar(64) UNIQUE NOT NULL,
    password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id serial PRIMARY KEY,
    user_id int NOT NULL,
    category varchar(64) NOT NULL,
    created timestamptz NOT NULL,
    score int NOT NULL,
    text text NOT NULL,
    title text NOT NULL,
    type varchar(64) NOT NULL,
    upvote_percentage numeric NOT NULL,
    url varchar(64) NOT NULL,
    views int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS votes (
    id serial PRIMARY KEY,
    user_id int NOT NULL,
	post_id int NOT NULL,
    vote int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    id serial PRIMARY KEY,
    user_id int NOT NULL,
	post_id int NOT NULL,
    body text NOT NULL,
    created timestamptz NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION score_upd() RETURNS trigger AS $$
    DECLARE
        count_upvotes int;
        count_downvotes int;
        recalc_upvote_percentage numeric;
    BEGIN
        IF TG_OP = 'DELETE' THEN
            SELECT SUM(vts.vote)
            FROM votes vts
            WHERE vts.post_id=OLD.post_id AND vts.vote=1
            INTO STRICT count_upvotes;
            IF count_upvotes IS NULL THEN
                count_upvotes := 0;
            END IF;

            SELECT SUM(vts.vote)
            FROM votes vts
            WHERE vts.post_id=OLD.post_id AND vts.vote=-1
            INTO STRICT count_downvotes;
            IF count_downvotes IS NULL THEN
                count_downvotes := 0;
            END IF;

            IF (count_upvotes + count_downvotes) = 0 THEN
                recalc_upvote_percentage := 0;
            ELSEIF count_upvotes = 0 THEN
                recalc_upvote_percentage := 0;
            ELSEIF count_downvotes = 0 THEN
                recalc_upvote_percentage := 100;
            ELSE
                recalc_upvote_percentage := (count_upvotes / (count_upvotes + count_downvotes));
            END IF;

            UPDATE posts
            SET (score, upvote_percentage) = (count_upvotes + count_downvotes, recalc_upvote_percentage)
            WHERE id=OLD.post_id;

            RETURN OLD;
        ELSEIF TG_OP = 'INSERT' THEN
            SELECT SUM(vts.vote)
            FROM votes vts
            WHERE vts.post_id=NEW.post_id AND vts.vote=1
            INTO STRICT count_upvotes;
            IF count_upvotes IS NULL THEN
                count_upvotes := 0;
            END IF;

            SELECT SUM(vts.vote)
            FROM votes vts
            WHERE vts.post_id=NEW.post_id AND vts.vote=-1
            INTO STRICT count_downvotes;
            IF count_downvotes IS NULL THEN
                count_downvotes := 0;
            END IF;

            IF (count_upvotes + count_downvotes) = 0 THEN
                recalc_upvote_percentage := 0;
            ELSEIF count_upvotes = 0 THEN
                recalc_upvote_percentage := 0;
            ELSEIF count_downvotes = 0 THEN
                recalc_upvote_percentage := 100;
            ELSE
                recalc_upvote_percentage := (count_upvotes / (count_upvotes + count_downvotes));
            END IF;

            UPDATE posts
            SET (score, upvote_percentage) = (count_upvotes + count_downvotes, recalc_upvote_percentage)
            WHERE id=NEW.post_id;

            RETURN NEW;
        END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER score_upd AFTER DELETE OR INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE score_upd();

INSERT INTO users (id, username, password) VALUES
(0, 'username', '65437890');

INSERT INTO posts (id, user_id, category, created, score, text, title, type, upvote_percentage, url, views) VALUES
(0, 0, 'music', '2014-04-04 20:00:00-07', 5, 'text content', 'title content', 'text', 15.5, '', 10);

INSERT INTO votes (id, user_id, post_id, vote) VALUES
(0, 0, 0, -1);

INSERT INTO comments (id, user_id, post_id, body, created) VALUES
(0, 0, 0, 'comment content', '2014-04-04 20:00:00-07');