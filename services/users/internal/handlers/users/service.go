package users

import (
	"screw-it-desk-backend/services/users/internal/service"
)

type Implementation struct {
	noteService service.NoteService
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
