package users

import (
	"database/sql"

	mysqldb "github.com/oscaralmgren/hackernews/internal/pkg/db/migrations/mysql"
	"golang.org/x/crypto/bcrypt"

	"log"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"name" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func (user *User) Create() {
	statement, err := mysqldb.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err2 := HashPassword(user.Password)
	if err2 != nil {
		log.Fatal(err2)
	}
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func (user *User) Authenticate() bool {
	statement, err := mysqldb.Db.Prepare("SELECT PASSWORD FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := mysqldb.Db.Prepare("SELECT ID FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}
