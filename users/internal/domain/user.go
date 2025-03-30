package domain

import "time"

type User struct {
	Id          string
	Username    string
	Email       string
	FullName    *string
	DateOfBirth *string
	PhoneNumber *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLogin   time.Time
	IsActive    bool
}
