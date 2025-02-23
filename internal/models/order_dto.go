package models

import "time"

type Order struct {
	ID       int       `db:"id"`
	PetID    int       `db:"petid"`
	Quantity int       `db:"quantity"`
	ShipDate time.Time `db:"shipdate"`
	Status   string    `db:"status"`
	Complete bool      `db:"complete"`
}
