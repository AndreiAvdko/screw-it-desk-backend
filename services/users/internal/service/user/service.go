package note

import (
	"screw-it-desk-backend/services/users/internal/client/db"
	"screw-it-desk-backend/services/users/internal/repository"
	"screw-it-desk-backend/services/users/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}
