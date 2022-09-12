package links

import (
	mysqldb "github.com/oscaralmgren/hackernews/internal/pkg/db/migrations/mysql"
	"github.com/oscaralmgren/hackernews/internal/users"

	"log"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

// #2
func (link Link) Save() int64 {
	//#3
	stmt, err := mysqldb.Db.Prepare("INSERT INTO Links(Title,Address) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(link.Title, link.Address)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted! Id:", id)
	return id
}

func GetAll() []Link {
	stmt, err := mysqldb.Db.Prepare("SELECT id, title, address FROM Links")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
