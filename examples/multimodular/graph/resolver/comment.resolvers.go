package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"multimodular/auth"
	storage1 "multimodular/comment/storage"
	"multimodular/post/storage"

	"github.com/gofrs/uuid"
	pgx "github.com/jackc/pgx/v5"
)

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, request storage1.DeleteCommentParams) (bool, error) {
	currentUserID := auth.GetCurrentUserId(ctx)
	if currentUserID == nil {
		return false, auth.ErrUnauthenticated
	}
	request.AuthorID = *currentUserID

	return true, r.CommentQueries.DeleteComment(ctx, request)
}

// LeaveComment is the resolver for the leaveComment field.
func (r *mutationResolver) LeaveComment(ctx context.Context, request storage1.LeaveCommentParams) (
	storage1.Comment,
	error,
) {
	currentUserID := auth.GetCurrentUserId(ctx)
	if currentUserID == nil {
		return storage1.Comment{}, auth.ErrUnauthenticated
	}
	request.ID = uuid.Must(uuid.NewV6())
	request.AuthorID = *currentUserID

	err := request.ValidateWithContext(ctx)
	if err != nil {
		return storage1.Comment{}, err
	}

	_, err = r.PostQueries.GetPost(ctx, request.PostID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage1.Comment{}, errors.New("post not found")
		}
		return storage1.Comment{}, err
	}

	return r.CommentQueries.LeaveComment(ctx, request)
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(
	ctx context.Context,
	obj *storage.Post,
	request storage1.GetPostCommentsParams,
) ([]storage1.Comment, error) {
	request.PostID = obj.ID
	return r.CommentQueries.GetPostComments(ctx, request)
}
