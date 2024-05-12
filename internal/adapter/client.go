package adapter

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"iman_tg_bot/internal/model"
)

type ClientI interface {
	CreateOne(res *model.Client) (*model.Client, error)
}

type clientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (c *clientRepo) CreateOne(req *model.Client) (*model.Client, error) {

	result := model.Client{}

	query := `
		INSERT INTO client_info (
		    contract_id,
		    phone_number,
		  	address,
		    payment_sum, 
		    comment, 
		    location,
			address_foto_path,
		    payment_foto_path,
		    chat_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING 
			contract_id,
		    phone_number,
		    address,
		    payment_sum, 
		    comment, 
		    location,
		    address_foto_path,
		    payment_foto_path,
		    chat_id
`
	if err := c.db.QueryRow(context.Background(), query,
		req.ContractId,
		req.PhoneNumber,
		req.Address,
		req.PaymentSum,
		req.Comment,
		req.Location,
		req.AddressFotoPath,
		req.PaymentFotoPath,
		req.ChatId,
	).Scan(
		&result.ContractId,
		&result.PhoneNumber,
		&result.Address,
		&result.PaymentSum,
		&result.Comment,
		&result.Location,
		&result.AddressFotoPath,
		&result.PaymentFotoPath,
		&result.ChatId,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
