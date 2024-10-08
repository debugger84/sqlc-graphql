package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"simple/graph/gen"
	"simple/storage"
)

// Bio is the resolver for the bio field.
func (r *authorResolver) Bio(ctx context.Context, obj *storage.Author) (*string, error) {
	if obj.Bio.Valid {
		return &obj.Bio.String, nil
	}
	return nil, nil
}

// Author returns gen.AuthorResolver implementation.
func (r *Resolver) Author() gen.AuthorResolver { return &authorResolver{r} }

type authorResolver struct{ *Resolver }
