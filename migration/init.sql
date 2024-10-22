CREATE TABLE IF NOT EXISTS categories (
    category_id INTEGER PRIMARY KEY,
    category_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    expire_time TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS posts (
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    likes_count INT NULL,
    dislikes_count INT NULL,
    comments_count INT NULL,
    tags TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post_likes(
    like_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS post_dislikes(
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
    likes_count INT NOT NULL,
    dislikes_count INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_likes(
    like_id TEXT NOT NULL,
    comment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_dislikes(
    dislike_id TEXT NOT NULL,
    comment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id)  ON DELETE CASCADE
);


INSERT INTO comments(comment_id, post_id, author, comment_text, likes_count, dislikes_count, created_at, updated_at)
VALUES('1', '', 'author', 'This is a comment\nwith a new line', 0, 0, datetime('now', 'localtime'), datetime('now', 'localtime'));

INSERT OR IGNORE INTO categories (category_name) 
VALUES ('Art'), ('Music'), ('Cinema'), ('Dance'), ('Architecture'), ('Fashion'), ('Graphics'), ('Cooking'),('IT'),('Life');



