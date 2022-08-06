package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
)

// PostgresDriver provides our implementation for the
// sql package.
type PostgresDriver struct{}

// Open provides a connection to the database.
func (dr PostgresDriver) Open(name string) (driver.Conn, error) {
	fmt.Println(name)
	return nil, errors.New("Unimplemented")
}

var d *PostgresDriver

// init is called prior to main.
func init() {
	fmt.Println("init")
	d = new(PostgresDriver)
	sql.Register("postgres", d)
}
