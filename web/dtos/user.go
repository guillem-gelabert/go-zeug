package dto

import "fmt"

type SignupDTO struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto SignupDTO) ValidateDisplayName(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto SignupDTO) ValidateEmail(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto SignupDTO) ValidatePassword(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto LoginDTO) ValidateEmail(predicate func(s string) bool, msg string) error {
	if !predicate(dto.Email) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto LoginDTO) ValidatePassword(predicate func(s string) bool, msg string) error {
	if !predicate(dto.Password) {
		return fmt.Errorf(msg)
	}
	return nil
}
