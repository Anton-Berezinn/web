package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type DataUser struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

var (
	EmptyDataError = errors.New("data is empty")
)

func DataRequest(data io.ReadCloser) (DataUser, error) {
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return DataUser{}, err
	}
	defer data.Close()
	var input DataUser
	err = json.Unmarshal(body, &input)
	if err != nil {
		return DataUser{}, err
	}
	if input.Lastname == "" || input.Login == "" || input.Password == "" || input.Firstname == "" {
		return DataUser{}, fmt.Errorf("%w", EmptyDataError)
	}
	return input, nil
}

func DataLogin(data io.ReadCloser) (*DataUser, error) {
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	input := &DataUser{}
	err = json.Unmarshal(body, input)
	if err != nil {
		return nil, err
	}
	if input.Login == "" || input.Password == "" {
		return nil, fmt.Errorf("%w", EmptyDataError)
	}
	return input, nil
}
