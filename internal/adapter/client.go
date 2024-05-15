package adapter

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"iman_tg_bot/internal/model"
	"log"
)

type ClientI interface {
	Create(chatId int64, username string) (*model.Client, error)
	Get(chatId int64) (*model.Client, error)
	GetOrCreate(chatId int64, username string) (*model.Client, error)
	UpdateOneFild(chatId int64, fild string, value string) error
	ChangeStep(ChatID int64, step float64) error
}

type clientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (r clientRepo) Create(chatId int64, username string) (*model.Client, error) {
	res := model.Client{}
	query := `
	INSERT INTO client_info(
		chat_id, user_name
	) VALUES(
		$1, $2
	) RETURNING chat_id, step`
	err := r.db.QueryRow(context.Background(), query, chatId, username).Scan(&res.ChatId, &res.Step)
	log.Println("Error", err)
	if err != nil {
		return &model.Client{}, err
	}
	return &res, err
}

func (r clientRepo) Get(ChatId int64) (*model.Client, error) {
	res := model.Client{}
	query := `
		SELECT 
		    contract_id,
    		phone_number,
   			address,
    		payment_sum,
    		comment,
    		location,
    		address_foto_path,
    		payment_foto_path,
    		user_name,
    		chat_id,
    		step
		FROM 
		    client_info 
		WHERE chat_id=$1`
	err := r.db.QueryRow(context.Background(), query, ChatId).Scan(
		&res.ContractId,
		&res.PhoneNumber,
		&res.Address,
		&res.PaymentSum,
		&res.Comment,
		&res.Location,
		&res.AddressFotoPath,
		&res.PaymentFotoPath,
		&res.UserName,
		&res.ChatId,
		&res.Step,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r clientRepo) GetOrCreate(chatId int64, username string) (*model.Client, error) {
	user, err := r.Get(chatId)
	if errors.Is(err, pgx.ErrNoRows) {
		u, err := r.Create(chatId, username)
		if err != nil {
			return &model.Client{}, err
		}
		user = u
	}
	return user, nil
}

func (r clientRepo) UpdateOneFild(chatId int64, fild string, value string) error {
	query := `UPDATE client_info SET ` + fild + `=$1 WHERE chat_id=$2`
	_, err := r.db.Exec(context.Background(), query, value, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (r clientRepo) ChangeStep(ChatID int64, step float64) error {
	query := `UPDATE client_info SET step=$1 WHERE chat_id=$2`
	_, err := r.db.Exec(context.Background(), query, step, ChatID)
	if err != nil {
		return err
	}
	return nil
}
