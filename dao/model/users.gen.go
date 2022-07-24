// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID              int64      `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	Name            string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email           string     `gorm:"column:email;type:varchar(255);not null" json:"email"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at;type:timestamp" json:"email_verified_at"`
	Password        string     `gorm:"column:password;type:varchar(255);not null" json:"password"`
	RememberToken   *string    `gorm:"column:remember_token;type:varchar(100)" json:"remember_token"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
