package register

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type StoragRegister struct {
	db *sql.DB
}

type Storage interface {
	Register(name, surname, login, password string) (bool, error)
}

var (
	ConnectDb        = errors.New("error to connect to database")
	DataError        = errors.New("error to add to database")
	Rowserror        = errors.New("rows is 0")
	CreateTableError = errors.New("error to create table")
)

func NewStorageRegister(data string) (*StoragRegister, error) {
	db, err := sql.Open("postgres", data)
	if err != nil {
		return nil, fmt.Errorf("%w", ConnectDb)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w", ConnectDb)
	}
	err = CreateTable(db)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err, CreateTableError)
	}
	return &StoragRegister{db: db}, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS users")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE DATABASE users(" +
		"id INT PRIMARY KEY," +
		"name VARCHAR(255)," +
		"surname VARCHAR(255)," +
		"login VARCHAR(255), " +
		"password VARCHAR(255), " +
		"CHECK (LENGTH(password) > 6)," +
		"CHECK (LENGTH(login) > 6);")
	if err != nil {
		return err
	}
	return nil
}

func (s *StoragRegister) Register(name, surname, login, password string) (bool, error) {
	count, err := s.db.Exec("INSERT INTO users(name,surname,login,password)VALUES($1,$2,$3,$4);", name, surname, login, password)
	if err != nil {
		return false, fmt.Errorf("%w", DataError)
	}
	value, err := count.RowsAffected()
	if err != nil || value != 1 {
		return false, fmt.Errorf("%w", Rowserror)
	}
	return true, nil
}
