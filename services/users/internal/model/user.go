package model

import (
	"database/sql"
	"time"
)

// TODO описать тип данных пользователя вместо шаблона для заметок

type Note struct {
	ID        int64
	Info      NoteInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type NoteInfo struct {
	Title   string
	Content string
}
