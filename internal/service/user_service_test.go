package service

import (
	"context"
	"testing"

	"go101/internal/model"
	"go101/internal/repository"
)

func TestCreateUser_ValidInput(t *testing.T) {
	service := &UserService{repo: nil}

	_, err := service.CreateUser(context.Background(), "Alice", "alice@example.com")
	if err == nil {
		t.Fatal("expected repository initialization error")
	}
}

func TestCreateUser_RejectsInvalidEmail(t *testing.T) {
	service := &UserService{repo: &repository.UserRepository{}}

	_, err := service.CreateUser(context.Background(), "Alice", "not-an-email")
	if err == nil {
		t.Fatal("expected invalid email error")
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	service := &UserService{repo: nil}

	_, err := service.GetUser(context.Background(), 0)
	if err == nil {
		t.Fatal("expected invalid user id error")
	}
}

func TestIsValidEmail(t *testing.T) {
	if !isValidEmail("demo@example.com") {
		t.Fatal("expected valid email")
	}
	if isValidEmail("bad-email") {
		t.Fatal("expected invalid email")
	}
}

func TestModelUser(t *testing.T) {
	user := model.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	if user.Name != "Alice" {
		t.Fatal("expected user name to be set")
	}
}
