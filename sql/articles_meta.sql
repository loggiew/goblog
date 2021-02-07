create table article_meta(
	   article_id INT NOT NULL AUTO_INCREMENT,
	   article_title VARCHAR(100) NOT NULL DEFAULT " ",
	   article_author VARCHAR(40) NOT NULL DEFAULT " ",
	   description VARCHAR(100) DEFAULT  " ",
	   submission_date DATE,
	   PRIMARY KEY ( article_id )
) ENGINE=InnoDB;
