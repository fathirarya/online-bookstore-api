package utils

import (
	"errors"
	"regexp"
)

var (
	reUpper      = regexp.MustCompile(`[A-Z]`)
	reLower      = regexp.MustCompile(`[a-z]`)
	reDigit      = regexp.MustCompile(`[0-9]`)
	reSymbol     = regexp.MustCompile(`[^A-Za-z0-9]`) // selain huruf/angka
	reWhitespace = regexp.MustCompile(`\s`)
)

// ValidatePassword mengecek apakah password memenuhi aturan:
// - Minimal 8 karakter
// - Ada huruf besar, huruf kecil, angka, dan simbol
// - Tidak boleh ada spasi/whitespace
func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !reUpper.MatchString(pw) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !reLower.MatchString(pw) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !reDigit.MatchString(pw) {
		return errors.New("password must contain at least one number")
	}
	if !reSymbol.MatchString(pw) {
		return errors.New("password must contain at least one symbol")
	}
	if reWhitespace.MatchString(pw) {
		return errors.New("password must not contain spaces")
	}
	return nil
}
