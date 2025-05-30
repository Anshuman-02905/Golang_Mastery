package dbrepo

import (
	"errors"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
)

func (*testDbRepo) AllUsers() bool {
	return true
}

// InsertReservaton inserts a reservation into the data
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return -1, errors.New("New Error")
	}
	return 1, nil
}

// inserts a room restriction in database
func (m *testDbRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 3 {
		return errors.New("New Error")
	}

	return nil

}

// SearchAvailabilityByDatesByRoomID Returns true if avaibiliryt exitst for a room
func (m *testDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	if roomID == 2 {
		return false, errors.New("New Error")
	}

	return true, nil

}

// SearchAvailabilityForAllRooms returns a aslice of available rooms  if any
func (m *testDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	cmpDate, _ := time.Parse("2006-01-02", "1990-04-30")
	if start.Equal(cmpDate) {
		return rooms, errors.New("New Error")
	}
	cmpDate, _ = time.Parse("2006-01-02", "1990-04-30")
	if end.Equal(cmpDate) {
		return rooms, nil
	}
	rooms = append(rooms, models.Room{ID: 1, RoomName: "General's Quaters", CreatedAt: time.Now(), UpdatedAt: time.Now()})

	return rooms, nil

}

// Gets a room by ID
func (m *testDbRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("Some Error")
	}

	return room, nil

}

func (m *testDbRepo) GetUserbyID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDbRepo) UpdatedUser(u models.User) error {
	return nil
}

func (m *testDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
