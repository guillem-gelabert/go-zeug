package dto

import "fmt"

type UserDTO struct {
	DisplayName string `display_name`
	Email       string `email`
	Password    string `password`
}

func (dto UserDTO) ValidateDisplayName(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto UserDTO) ValidateEmail(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}

func (dto UserDTO) ValidatePassword(predicate func(s string) bool, msg string) error {
	if !predicate(dto.DisplayName) {
		return fmt.Errorf(msg)
	}
	return nil
}
