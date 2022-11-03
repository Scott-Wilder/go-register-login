Basic go lang account registration and login with DB:

- User registration html form with authentication.
- DB connection with mysql.
- Read db password from OS env variables to protect password.
- Encrypt user password and store only hash in db to protect password. (bcrpyt to encrypt passoword)

- DB Creation: CREATE TABLE User (userID int NOT NULL AUTO_INCREMENT, username lastname varchar(25) NOT NULL, firstname varchar(25) NOT NULL, Hash varchar(255), PRIMARY KEY (userID));
