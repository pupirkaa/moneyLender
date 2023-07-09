package ml

import (
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Users UsersStorage
}

var (
	ErrUserNotFound    = errors.New("can't find a user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidSignup   = errors.New("invalid name or password")
)

func (s *AuthService) Login(name string, password string) (session string, err error) {
	exists, err := s.Users.UserExist(name)
	if err != nil {
		return "", fmt.Errorf("getting user:%v", err)
	}

	if !exists {
		return "", ErrUserNotFound
	}

	hashedPassword, err := s.Users.UserGet(name)
	if err != nil {
		return "", fmt.Errorf("getting user's password:%v", err)
	}

	if !ComparePasswords(hashedPassword, password) {
		return "", ErrInvalidPassword
	}

	return MakeSession(name), nil
}

func (s *AuthService) Signup(name string, password string) (err error) {
	if !IsNameValid(name) || !IsPasswordValid(password) {
		return ErrInvalidSignup
	}
	hashedPassword, err := HashAndSalt(password)
	if err != nil {
		return fmt.Errorf("hashing user's password:%v", err)
	}

	s.Users.UserAdd(name, hashedPassword)
	return nil
}

func IsNameValid(name string) bool {
	return !regexp.MustCompile(`\s`).MatchString(name)
}

func IsPasswordValid(password string) bool {
	return (!regexp.MustCompile(`\s`).MatchString(password) || len(password) < 4)
}

func HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %v", err)
	}
	return string(hash), nil
}

func MakeSession(name string) string {
	return hex.EncodeToString([]byte(name + time.Now().GoString()))
}

func ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		fmt.Println("comparing passwords: ", err)
		return false
	}
	return true
}
