package handlers

import (
	"screw-it-desk-backend/services/users/internal/service"
)

type Implementation struct {
	noteService service.UserService
}

func NewImplementation(noteService service.UserService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
