create table content(
	   content_id INT NOT NULL AUTO_INCREMENT,
	   content TEXT NOT NULL,
	   last_updated_date DATE,
	   PRIMARY KEY (content_id)
) ENGINE=InnoDB;
