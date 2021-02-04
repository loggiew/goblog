INSERT INTO article_meta (article_title, article_author) VALUES ("Solar Panels", "Logan");
INSERT INTO article_meta (article_title, article_author) VALUES ("Drones", "Bob");
INSERT INTO article_meta (article_title, article_author) VALUES ("3d Dev", "Fred");
INSERT INTO article_meta (article_title, article_author) VALUES ("Blender", "Mike");
INSERT INTO article_meta (article_title, article_author) VALUES ("Sewing", "Jim");
INSERT INTO article_meta (article_title, article_author) VALUES ("Cooking", "Edward");

INSERT INTO content (content) VALUES ("This is my article. There are many like it, but this article is mine.");
INSERT INTO content (content) VALUES ("");
INSERT INTO content (content) VALUES ("");
INSERT INTO content (content) VALUES ("");
INSERT INTO content (content) VALUES ("");
INSERT INTO content (content) VALUES ("");

INSERT INTO articles (content_id,article_id) VALUES (1,1);
INSERT INTO articles (content_id,article_id) VALUES (2,2);
INSERT INTO articles (content_id,article_id) VALUES (3,3);
INSERT INTO articles (content_id,article_id) VALUES (4,4);
INSERT INTO articles (content_id,article_id) VALUES (5,5);
INSERT INTO articles (content_id,article_id) VALUES (6,6);

#drop table article_meta;
#drop table content;
#drop table articles;

