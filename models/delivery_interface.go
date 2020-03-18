package models

type DeliveryTranslator interface {
	Translate() (Delivery, error)
}
