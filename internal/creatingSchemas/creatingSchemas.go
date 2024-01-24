package creatingSchemas

import (
	"github.com/go-pg/pg/orm"
	"l0/internal/database"
	"l0/internal/model"
)

func CreateDbSchemas(db *database.Database) error {
	orders := []interface{}{(*model.Order)(nil)}
	err := CreateSchema(db, orders)
	if err != nil {
		return err
	}

	deliveries := []interface{}{(*model.Delivery)(nil)}
	err = CreateSchema(db, deliveries)
	if err != nil {
		return err
	}

	payments := []interface{}{(*model.Payment)(nil)}
	err = CreateSchema(db, payments)
	if err != nil {
		return err
	}

	items := []interface{}{(*model.Item)(nil)}
	err = CreateSchema(db, items)
	if err != nil {
		return err
	}

	return nil
}

func CreateSchema(db *database.Database, models []interface{}) error {
	for _, mod := range models {
		op := orm.CreateTableOptions{}
		err := db.DB.Model(mod).CreateTable(&op)
		if err != nil {
			return err
		}
	}
	return nil
}
