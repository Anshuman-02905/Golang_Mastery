package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/jackc/pgconn"
)

// DB holds the database connection Pool
type DB struct {
	SQL *sql.DB
}

var dbCOnn = &DB{}

const MaxOpenDbConn = 10
const MaxIdleDbConn = 5
const DbLifetime = 5 * time.Minute

// Connect SQL creates Database pools for POstgress
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(MaxOpenDbConn)
	d.SetConnMaxLifetime(DbLifetime)
	d.SetMaxIdleConns(MaxIdleDbConn) //WONT GROW OUT OF CONTROL

	dbCOnn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbCOnn, nil
}

// Tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
