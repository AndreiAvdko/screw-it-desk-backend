package model

import (
	"database/sql"
	"time"
)

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

// TODO
// CREATE TABLE users (
//     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     email VARCHAR(255) NOT NULL UNIQUE,
//     username VARCHAR(100) NOT NULL UNIQUE,
//     password_hash VARCHAR(255) NOT NULL,

//     -- Личная информация
//     first_name VARCHAR(100),
//     last_name VARCHAR(100),
//     phone VARCHAR(20),

//     -- Статусы и роли
//     role VARCHAR(50) DEFAULT 'user',
//     status VARCHAR(50) DEFAULT 'active',

//     -- Метаданные
//     email_verified_at TIMESTAMP,
//     phone_verified_at TIMESTAMP,
//     last_login_at TIMESTAMP,
//     last_login_ip INET,

//     -- Временные метки
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     deleted_at TIMESTAMP, -- для soft delete
