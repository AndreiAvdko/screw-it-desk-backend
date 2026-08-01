package note

import (
	"screw-it-desk-backend/services/users/internal/client/db"
	"screw-it-desk-backend/services/users/internal/repository"
	"screw-it-desk-backend/services/users/internal/service"
)

type serv struct {
	noteRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}
