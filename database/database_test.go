package database

import "testing"

func TestDb(t *testing.T) {
	c, err := NewClient()

	if err != nil {
		t.Errorf("couldn't create a db client: %s", err)
		return
	}

	if err = c.Connection.Ping(); err != nil {
		t.Errorf("couldn't ping Postgres: %s", err)
		return
	}

	var id int
	if err = c.Connection.QueryRow(
		`INSERT INTO brands(name) VALUES ("sainsbury's") returning id`,
	).Scan(&id); err != nil {
		t.Errorf("couldn't insert data: %s", err)
		return
	}

	if id == 0 {
		t.Error("the inserted id doesn't look right")
	}

	// We don't need to test the underlying behavior of the database client
	// as it is a dependency, which is tested in its own repository.
}
