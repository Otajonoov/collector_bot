package adapter

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	ClientUser() ClientI
}

type repo struct {
	clientRepo ClientI
}

func NewRepo(db *pgxpool.Pool) *repo {
	return &repo{
		clientRepo: NewClientRepo(db),
	}
}

func (r *repo) ClientUser() ClientI {
	return r.clientRepo
}
