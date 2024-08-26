// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"database/sql/driver"
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusDeleted   PostStatus = "deleted"
)

func (e *PostStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PostStatus(s)
	case string:
		*e = PostStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PostStatus: %T", src)
	}
	return nil
}

type NullPostStatus struct {
	PostStatus PostStatus `json:"postStatus"`
	Valid      bool       `json:"valid"` // Valid is true if PostStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPostStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PostStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PostStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPostStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PostStatus), nil
}

func AllPostStatusValues() []PostStatus {
	return []PostStatus{
		PostStatusDraft,
		PostStatusPublished,
		PostStatusDeleted,
	}
}

type Post struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Status      PostStatus         `json:"status"`
	AuthorID    uuid.UUID          `json:"authorId"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	PublishedAt pgtype.Timestamptz `json:"publishedAt"`
}
