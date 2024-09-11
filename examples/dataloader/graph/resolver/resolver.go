package resolver

import (
	"dataloader/storage"
	"dataloader/storage/dataloader"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Queries       *storage.Queries
	LoaderFactory *dataloader.LoaderFactory
}
