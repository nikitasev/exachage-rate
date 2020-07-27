package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type DbService struct {
	Con    *sql.DB
	Locker sync.Mutex
}

func InitDBService() (dbService DbService) {
	db, err := sql.Open("mysql", "root@(127.0.0.1:3306)/test")
	dbService = DbService{Con: db}

	if err != nil {
		panic(err.Error())
	}

	return dbService
}
