package model

type Client struct {
	Id                int
	ContractId        string
	PhoneNumber       string
	Address           string
	PaymentSum        string
	Comment           string
	Location          string
	LocationLatitude  string
	LocationLongitude string
	AddressFoto       string
	PaymentFoto       string
}

type UserPass struct {
	Password string
}
