package dataloader

import (
	"dataloader/storage"
)

type LoaderFactory struct {
	db           storage.DBTX
	authorLoader *AuthorLoader
	postLoader   *PostLoader
}

func NewLoaderFactory(db storage.DBTX) *LoaderFactory {
	return &LoaderFactory{
		db: db,
	}
}

func (f *LoaderFactory) AuthorLoader() *AuthorLoader {
	if f.authorLoader == nil {
		f.authorLoader = NewAuthorLoader(f.db)
	}
	return f.authorLoader
}
func (f *LoaderFactory) PostLoader() *PostLoader {
	if f.postLoader == nil {
		f.postLoader = NewPostLoader(f.db)
	}
	return f.postLoader
}
