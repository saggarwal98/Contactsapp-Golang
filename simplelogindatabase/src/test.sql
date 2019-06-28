Use simplelogin;

Create table `Users` (ID int PRIMARY KEY AUTO_INCREMENT,FIRSTNAME varchar(255),LASTNAME varchar(255),EMAIL varchar(40) NOT NULL,PASSWORD varchar(50) NOT NULL);

Insert into Users (FIRSTNAME,LASTNAME,EMAIL,PASSWORD) values ("anonymous","user","abc@xyz.com","anonymous");