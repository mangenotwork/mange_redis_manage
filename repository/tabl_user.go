package models

/*
CREATE TABLE table_user(
   uid INTEGER PRIMARY KEY,
   uname TEXT NOT NULL,
   upassword TEXT NOT NULL,
   ugroup INT NOT NULL
);
*/

type User struct {
	Uid       int64
	Uname     string
	Upassword string
	Ugroup    int64
}
