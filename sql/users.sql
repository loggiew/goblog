create table users(
	   user_id INT NOT NULL AUTO_INCREMENT,
	   username VARCHAR(255) NOT NULL,
	   password VARCHAR(255) NOT NULL,
	   last_login_date DATE,
	   PRIMARY KEY (user_id)
) ENGINE=InnoDB;
