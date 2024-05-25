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
	UserName        string
	ChatId          int64
	Step            float64
}

const (
	StartCommandStep  = 1
	CheckUserPassword = 2
	AddData           = 3
	ContractId        = 4
	PhoneNumber       = 4.1
	Address           = 4.2
	PaymentSum        = 4.3
	Comment           = 4.4
	Location          = 4.5
	AddressFotoPath   = 4.6
	PaymentFotoPath   = 4.7
	Finish            = 5
)

const (
	MenuText = "<b>Assalomu alaykum.</b>"
)
