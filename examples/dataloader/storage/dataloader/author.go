package dataloader

import (
	"context"
	"dataloader/storage"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/jackc/pgx/v5"
)

type AuthorLoader struct {
	innerLoader *dataloader.Loader[int64, storage.Author]
	db          storage.DBTX
}

func NewAuthorLoader(db storage.DBTX) *AuthorLoader {
	return &AuthorLoader{
		db: db,
	}
}

func (l *AuthorLoader) getInnerLoader() *dataloader.Loader[int64, storage.Author] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []int64) []*dataloader.Result[storage.Author] {
				authorMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Author], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Author]{Error: err}
						continue
					}

					if loadedItem, ok := authorMap[key]; ok {
						result[i] = &dataloader.Result[storage.Author]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Author]{Error: pgx.ErrNoRows}
					}
				}
				return result
			},
		)
	}
	return l.innerLoader
}

func (l *AuthorLoader) findItemsMap(ctx context.Context, keys []int64) (map[int64]storage.Author, error) {
	res := make(map[int64]storage.Author, len(keys))

	query := `SELECT * FROM authors WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Author
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Bio,
			&result.Status,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *AuthorLoader) Load(ctx context.Context, authorKey int64) (storage.Author, error) {
	return l.getInnerLoader().Load(ctx, authorKey)()
}
