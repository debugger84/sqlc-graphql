package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"multimodular/auth"
	"multimodular/post/storage"

	"github.com/gofrs/uuid"
)

// MakeDraft is the resolver for the makeDraft field.
func (r *mutationResolver) MakeDraft(ctx context.Context, request storage.MakeDraftParams) (storage.Post, error) {
	curentUserId := auth.GetCurrentUserId(ctx)
	if curentUserId == nil {
		return storage.Post{}, auth.ErrUnauthenticated
	}
	request.AuthorID = *curentUserId
	request.ID = uuid.Must(uuid.NewV6())
	return r.PostQueries.MakeDraft(ctx, request)
}

// Publish is the resolver for the publish field.
func (r *mutationResolver) Publish(ctx context.Context, id uuid.UUID) (storage.Post, error) {
	curentUserId := auth.GetCurrentUserId(ctx)
	if curentUserId == nil {
		return storage.Post{}, auth.ErrUnauthenticated
	}
	post, err := r.PostQueries.GetPost(ctx, id)
	if err != nil {
		return storage.Post{}, errors.New("post not found")
	}
	if post.AuthorID != *curentUserId {
		return storage.Post{}, errors.New("not the author")
	}

	return r.PostQueries.Publish(ctx, id)
}

// LastPosts is the resolver for the lastPosts field.
func (r *queryResolver) LastPosts(ctx context.Context, request storage.GetLastPostsParams) ([]storage.Post, error) {
	return r.PostQueries.GetLastPosts(ctx, request)
}

// MyDrafts is the resolver for the myDrafts field.
func (r *queryResolver) MyDrafts(ctx context.Context, request storage.GetMyDraftsParams) ([]storage.Post, error) {
	curentUserId := auth.GetCurrentUserId(ctx)
	if curentUserId == nil {
		return nil, auth.ErrUnauthenticated
	}

	request.AuthorID = *curentUserId
	return r.PostQueries.GetMyDrafts(ctx, request)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id uuid.UUID) (storage.Post, error) {
	return r.PostQueries.GetPost(ctx, id)
}
