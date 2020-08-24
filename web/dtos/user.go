package dto

import "fmt"

type SignupDTO struct {
	DisplayName string `display_name`
	Email       string `email`
	Password    string `password`
}

func (dto UserDTO) ValidateDisplayName(predicate func(s string) bool, msg string) error {

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
