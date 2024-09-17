CREATE TABLE IF NOT EXISTS categories (
    category_id INTEGER PRIMARY KEY,
    category_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    ID TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS session (
    session_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    expire_time TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS posts (
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    likes INT NULL,
    dislikes INT NULL,
    comments INT NULL,
    tags TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS likes(
    like_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS dislikes(
    dislike_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS comments (
    comment_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    author TEXT NOT NULL,
    comment_text TEXT NOT NULL,
    likes INT NOT NULL,
    dislikes INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_like(
    like_id TEXT NOT NULL,
    comment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_dislike(
    dislike_id TEXT NOT NULL,
    comment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id)  ON DELETE CASCADE
);


INSERT INTO comments(comment_id, post_id, author, comment_text, likes, dislikes, created_at, updated_at)
VALUES('1', '', 'author', 'This is a comment\nwith a new line', 0, 0, datetime('now', 'localtime'), datetime('now', 'localtime'));

INSERT OR IGNORE INTO categories (category_name) 
VALUES ('Art'), ('Music'), ('Cinema'), ('Dance'), ('Architecture'), ('Fashion'), ('Graphics'), ('Cooking'),('IT'),('Life');



