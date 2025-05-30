package models

import (
	"time"
)

// Users is the user Model
type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AcessLevel string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Rooms is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations is the Reservations Model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDAte   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestrictions is the RoomRestrictions model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDAte       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int

	CreatedAt    time.Time
	UpdatedAt    time.Time
	Room         Room
	Reservation  Reservation
	Restrictions Restriction
}

// mailData Holds mail data
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
