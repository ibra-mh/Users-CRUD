package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectDB() (x *sql.DB, err error){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		fmt.Println("==============================))))))))))))))))))))))))))))==========================================")
		return nil, err
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")

	if err != nil {
		log.Fatal(err)
		fmt.Print("___________________________________________________________________________")
		return nil, err
	}
	fmt.Println("**************************************************************************")
	return db, nil
}
