package repository

import (
	"context"

	"screw-it-desk-backend/services/users/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
