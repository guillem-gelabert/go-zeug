package dto

import "regexp"

func Required() func(s string) bool {
	return func(s string) bool {
		if s == "" {
			return false
		}
		return true
	}
}

func IsEmail() func(s string) bool {
	return func(s string) bool {
		emailRX := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
		if !emailRX.MatchString(s) {
			return false
		}
		return true
	}
}

func MaxLength(max int) func(s string) bool {
	return func(s string) bool {
		return len(s) <= max
	}
}

func MinLength(min int) func(s string) bool {
	return func(s string) bool {
		return len(s) >= min
	}
}
