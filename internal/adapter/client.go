package adapter

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"iman_tg_bot/internal/model"
)

type ClientI interface {
	Create(req *model.Client) error
}

type clientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (c *clientRepo) Create(req *model.Client) error {

	query := `
			INSERT INTO client_info (
			                         contract_id,
			                         phone_number,
			                         address,
			                         payment_sum,
			                         comment,
			                         location_latitude,
			                         location_longitude,
			                         address_foto_path,
			                         payment_foto_path
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.db.Exec(context.Background(), query,
		req.ContractId,
		req.PhoneNumber,
		req.Address,
		req.PaymentSum,
		req.Comment,
		req.LocationLatitude,
		req.LocationLongitude,
		req.AddressFoto,
		req.PaymentFoto,
	)
	if err != nil {
		return err
	}

	return nil
}
