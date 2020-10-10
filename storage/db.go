package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Storage is the db connection
var Storage *Db

// Db is the SQLite3 db handler ype
type Db struct {
	Conn *gorm.DB
}

// InitDb sets up a new DB
func InitDb() error {

	conn, err := gorm.Open(sqlite.Open("db.sqlite?cache=shared"), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Error),
		// todo: silence this when done
	})
	if err != nil {
		return err
	}

	Storage = &Db{Conn: conn}
	Storage.migrate()

	return nil
}

func (db *Db) migrate() {
	db.Conn.AutoMigrate(&Player{}, &Board{}, &Event{})
}

// Get gets a db handle
func (db *Db) Get() *gorm.DB {
	return db.Conn
}
