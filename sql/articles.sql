create table articles(
    content_id INT NOT NULL,
    article_id INT NOT NULL,
    PRIMARY KEY (content_id, article_id),
    FOREIGN KEY (content_id) references content (content_id) on delete cascade,
    FOREIGN KEY (article_id) references article_meta (article_id) on delete cascade
) ENGINE=InnoDB;
