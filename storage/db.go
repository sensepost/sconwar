package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Storage is the db connection
var Storage *Db

// Db is the SQLite3 db handler ype
type Db struct {
	Conn *gorm.DB
}

// InitDb sets up a new DB
func InitDb() error {
	return InitDbPath("db.sqlite")
}

// InitDbPath sets up a new DB at a specific sqlite file path.
func InitDbPath(path string) error {

	conn, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?cache=shared", path)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return err
	}

	Storage = &Db{Conn: conn}
	Storage.migrate()

	return nil
}

func (db *Db) migrate() {
	db.Conn.AutoMigrate(&Player{}, &Board{}, &Event{}, &PlayerGameScore{})
}

// Get gets a db handle
func (db *Db) Get() *gorm.DB {
	return db.Conn
}
