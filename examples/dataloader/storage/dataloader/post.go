package dataloader

import (
	"context"
	"dataloader/storage"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/jackc/pgx/v5"
)

type PostLoader struct {
	innerLoader *dataloader.Loader[int64, storage.Post]
	db          storage.DBTX
}

func NewPostLoader(db storage.DBTX) *PostLoader {
	return &PostLoader{
		db: db,
	}
}

func (l *PostLoader) getInnerLoader() *dataloader.Loader[int64, storage.Post] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []int64) []*dataloader.Result[storage.Post] {
				postMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Post], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Post]{Error: err}
						continue
					}

					if loadedItem, ok := postMap[key]; ok {
						result[i] = &dataloader.Result[storage.Post]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Post]{Error: pgx.ErrNoRows}
					}
				}
				return result
			},
		)
	}
	return l.innerLoader
}

func (l *PostLoader) findItemsMap(ctx context.Context, keys []int64) (map[int64]storage.Post, error) {
	res := make(map[int64]storage.Post, len(keys))

	query := `SELECT * FROM posts WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Post
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Content,
			&result.AuthorID,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *PostLoader) Load(ctx context.Context, postKey int64) (storage.Post, error) {
	return l.getInnerLoader().Load(ctx, postKey)()
}
