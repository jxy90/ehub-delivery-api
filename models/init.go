package models

import "github.com/go-xorm/xorm"

func Init(db *xorm.Engine) error {
	if err := db.Sync(
		new(Delivery), new(DeliveryItem)); err != nil {
		return err
	}
	return nil
}
