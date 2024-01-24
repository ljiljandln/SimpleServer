package database

import (
	"github.com/go-pg/pg"
	"l0/internal/model"
)

type config struct {
	Database string
	User     string
	Password string
}

type Database struct {
	DB     *pg.DB
	config *config
}

func SetConfig() *Database {
	config := config{Database: "l0_db", User: "postgres", Password: "test"}
	return &Database{nil, &config}
}

func (db *Database) Open() {
	db.DB = pg.Connect(&pg.Options{
		User:     db.config.User,
		Password: db.config.Password,
		Database: db.config.Database,
	})
}

func (db *Database) Close() {
	err := db.DB.Close()
	if err != nil {
		return
	}
}

func (db *Database) AddOrder(order model.Order) error {
	if _, err := db.DB.Model(&order).Insert(); err != nil {
		return err
	}
	return nil
}
