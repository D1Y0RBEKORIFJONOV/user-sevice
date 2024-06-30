package user_service

import (
	"log/slog"
	"time"
)

type User struct {
	log      *slog.Logger
	tokenTTL time.Duration
	provider UserProvider
	saver    UserSaver
	deleter  UserDeleter
	updater  UserUpdater
}

func NewUser(log *slog.Logger,
	provider UserProvider,
	saver UserSaver,
	deleter UserDeleter,
	updater UserUpdater,
	tokenTTl time.Duration) *User {
	return &User{
		log:      log,
		provider: provider,
		saver:    saver,
		deleter:  deleter,
		updater:  updater,
		tokenTTL: tokenTTl,
	}
}
