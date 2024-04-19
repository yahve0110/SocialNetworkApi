CREATE TABLE post_permissions (
    user_id TEXT,
    post_id TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    PRIMARY KEY (user_id, post_id)
);
