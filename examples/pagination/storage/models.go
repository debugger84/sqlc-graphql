// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuthorStatus string

const (
	AuthorStatusActive   AuthorStatus = "active"
	AuthorStatusInactive AuthorStatus = "inactive"
	AuthorStatusDeleted  AuthorStatus = "deleted"
)

func (e *AuthorStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AuthorStatus(s)
	case string:
		*e = AuthorStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AuthorStatus: %T", src)
	}
	return nil
}

type NullAuthorStatus struct {
	AuthorStatus AuthorStatus `json:"authorStatus"`
	Valid        bool         `json:"valid"` // Valid is true if AuthorStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAuthorStatus) Scan(value interface{}) error {
	if value == nil {
		ns.AuthorStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AuthorStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAuthorStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AuthorStatus), nil
}

func AllAuthorStatusValues() []AuthorStatus {
	return []AuthorStatus{
		AuthorStatusActive,
		AuthorStatusInactive,
		AuthorStatusDeleted,
	}
}

type Author struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Bio       pgtype.Text        `json:"bio"`
	Status    AuthorStatus       `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type AuthorPage struct {
	Items   []Author
	Total   int
	HasNext bool
}
