CREATE TABLE IF NOT EXISTS file_tags(
    file_id char(20) REFERENCES files(id) ON DELETE CASCADE NOT NULL,
    tag_id char(20) REFERENCES tags(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY(file_id, tag_id)
);
