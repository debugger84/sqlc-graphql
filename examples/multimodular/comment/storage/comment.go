package storage

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (p *LeaveCommentParams) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStruct(
		p,
		validation.Field(
			&p.Comment,
			validation.Required.ErrorObject(
				validation.NewError(
					"CommentRequired",
					"Comment is required",
				),
			),
			validation.Length(1, 1000).ErrorObject(
				validation.NewError(
					"CommentLength",
					"Comment length should be between 1 and 1000",
				),
			),
		),
		validation.Field(
			&p.PostID,
			validation.Required.ErrorObject(
				validation.NewError(
					"PostIDRequired",
					"PostID is required",
				),
			),
		),
	)
}
