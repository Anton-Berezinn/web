package register

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
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
	UpdateError      = errors.New("error to update table")
	RowsError        = errors.New("rows is 0")
	DeleteError      = errors.New("error to delete table")
	LenPassError     = errors.New("password len")
	NoLoginError     = errors.New("no login")
)

func NewStorageRegister(data string) (*StoragRegister, error) {
	db, err := sql.Open("postgres", data)
	if err != nil {
		return nil, fmt.Errorf("%w", ConnectDb)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w", ConnectDb)
	}
	return &StoragRegister{db: db}, nil
}

func (s *StoragRegister) Register(firstname, lastname, login, password string) error {
	count, err := s.db.Exec("INSERT INTO users(login,firstname,lastname,password)VALUES($1,$2,$3,$4);", login, firstname, lastname, password)
	if err != nil {
		fmt.Println("error to Exec error:", err)
		if strings.Contains(err.Error(), "CHECK (LENGTH(password) > 6)") {
			return fmt.Errorf("%w", LenPassError)
		}
		return fmt.Errorf("%w", DataError)
	}
	value, err := count.RowsAffected()
	if err != nil || value != 1 {
		return fmt.Errorf("%w", Rowserror)
	}
	return nil
}

func (s *StoragRegister) Update(login, password string) error {
	val, err := s.db.Exec("UPDATE users SET password=$1 WHERE login=$2;", password, login)
	if err != nil {
		return fmt.Errorf("%w", UpdateError)
	}
	value, err := val.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w", UpdateError)
	}
	if value != 1 {
		return fmt.Errorf("%w", RowsError)
	}
	return nil
}

func (s *StoragRegister) Delete(login string) error {
	val, err := s.db.Exec("DELETE FROM users WHERE login=$1;", login)
	if err != nil {
		return fmt.Errorf("%w", DeleteError)
	}
	value, err := val.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w", DeleteError)
	}
	if value != 1 {
		return fmt.Errorf("%w", RowsError)
	}
	return nil
}

func (s *StoragRegister) CheckPassword(login, passwordUser string) error {
	value, err := s.db.Query("SELECT password FROM users WHERE login=$1;", login)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return DataError
	}
	defer value.Close()
	if !value.Next() {
		return fmt.Errorf("%w", NoLoginError)
	}
	var password sql.NullString
	if err := value.Scan(&password); err != nil {
		return err
	}
	if !password.Valid || password.String != passwordUser {
		return fmt.Errorf("%w", NoLoginError)
	}
	return nil
}
