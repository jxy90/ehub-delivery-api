package models

type LocationType struct {
	Code string `json:"code"`
}

var (
	Store = LocationType{Code: "store"}
	Plant = LocationType{Code: "plant"}
)
