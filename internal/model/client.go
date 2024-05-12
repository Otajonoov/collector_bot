package model

type Client struct {
	Id              int
	ContractId      string
	PhoneNumber     string
	Address         string
	PaymentSum      string
	Comment         string
	Location        string
	AddressFotoPath string
	PaymentFotoPath string
	ChatId          int64
}
