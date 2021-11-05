package database

import "github.com/tarantool/go-tarantool"

type DB struct {
	DB *tarantool.Connection
}

func New(addr, user, pass string) (*DB, error) {
	db, err := tarantool.Connect(addr, tarantool.Opts{
		User:          user,
		Pass:          pass,
	})

	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}