package model

import "time"

type User struct {
	TelegramID        int64
	FirstName         string
	LastName          *string
	Username          *string
	SubscriptionToken *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UpsertUserParams struct {
	TelegramID int64
	FirstName  string
	LastName   *string
	Username   *string
}

type Subscription struct {
	ID             int64
	TelegramUserID int64
	PlanCode       string
	PlanLabel      string
	Status         string
	StartsAt       time.Time
	EndsAt         time.Time
	ClientUUID     *string
	IsMock         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateSubscriptionParams struct {
	TelegramUserID int64
	PlanCode       string
	PlanLabel      string
	StartsAt       time.Time
	EndsAt         time.Time
	ClientUUID     *string
	IsMock         bool
}
