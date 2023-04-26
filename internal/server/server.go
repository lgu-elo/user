package server

import (
	"github.com/lgu-elo/user/internal/user"
)

type (
	UserHandler user.IHandler

	API struct {
		UserHandler
	}
)

func NewAPI(user user.IHandler) *API {
	return &API{user}
}
