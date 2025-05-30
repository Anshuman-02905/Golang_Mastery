package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (*postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservaton inserts a reservation into the data
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,end_date,room_id,created_at,updated_at)
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDAte,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, err
}

// inserts a room restriction in database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restriction ( start_date,end_date,room_id,reservation_id,restriction_id,created_at,updated_at)
				values ($1,$2,$3,$4,$5,$6,$7) returning id`

	_, err := m.DB.ExecContext(ctx, stmt,

		res.StartDate,
		res.EndDAte,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil

}

// SearchAvailabilityByDatesByRoomID Returns true if avaibiliryt exitst for a room
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	// stmt := `
	// SELECT
	// 	COUNT(*)
	// FROM
	// 	room_restriction
	// WHERE room_id =$1
	// 	(start_date
	// BETWEEN
	// 	'$2' AND '$3' )
	// OR
	// 	(end_date
	// BETWEEN
	// 	'$2'  AND '$3')
	// OR (start_date<'$2' AND end_date>'$3');`
	stmt := `
	SELECT 
	COUNT(*) 
FROM 
	room_restriction 
WHERE 
	room_id = $1 AND (
		start_date BETWEEN $2 AND $3
		OR end_date BETWEEN $2 AND $3
		OR (start_date < $2 AND end_date > $3)
	);
	`
	row := m.DB.QueryRowContext(ctx, stmt, roomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}

	return false, nil

}

// SearchAvailabilityForAllRooms returns a aslice of available rooms  if any
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `SELECT 
    r.id, r.room_name
FROM
    rooms r
LEFT JOIN
    room_restriction rr 
    ON r.id = rr.room_id 
    AND (
        rr.start_date BETWEEN $1 AND $2
        OR rr.end_date BETWEEN $1 AND $2
        OR (rr.start_date < $1 AND rr.end_date > $2)
    )
WHERE rr.id IS NULL;`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID, &room.RoomName,
		)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil

}

// Gets a room by ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `
		SELECT id,room_name,created_at,updated_at from rooms WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil

}

// Returns a user BY ID
func (m *postgresDBRepo) GetUserbyID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email ,password, access_level , created_at, updated_at
			from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AcessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}
	return u, nil

}

//Update User in database

func (m *postgresDBRepo) UpdatedUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update user set first_name = $1 , last_name = $2,email = $3, access_level = $4,updated_at = $5,
	`
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AcessLevel,
		time.Now,
	)
	if err != nil {
		return err
	}
	return nil
}
//Authenticate authenticates the user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id,password from users where email = $1", email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect Password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
