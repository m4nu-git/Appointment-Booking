package domain

import "time"

type Role string

const (
	RoleUser             Role = "USER"
	RoleServiceProvider   Role = "SERVICE_PROVIDER"
)

// User maps onto the `users` table created in migration 000001.
// Table/column names match GORM's default snake_case convention, so no
// explicit `gorm:"column:..."` tags are needed here — but we're still
// listing the mapping intent for clarity.
type User struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	Role         Role      `gorm:"column:role"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (User) TableName() string { return "users" }