package storages

import (
	"pet-store/internal/db/adapter"
	ostorage "pet-store/internal/modules/order/storage"
	petstorage "pet-store/internal/modules/pet/storage"
	ustorage "pet-store/internal/modules/user/storage"
)

type Storages struct {
	User  ustorage.Userer
	Pet   petstorage.Peter
	Order ostorage.Orderer
}

func NewStorages(sqlAdapter *adapter.SQLAdapter) *Storages {
	return &Storages{
		User:  ustorage.NewUserStorage(sqlAdapter),
		Pet:   petstorage.NewPetStorage(sqlAdapter),
		Order: ostorage.NewOrderStorage(sqlAdapter),
	}
}
