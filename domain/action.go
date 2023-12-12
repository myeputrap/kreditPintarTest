package domain

import (
	"context"
	"time"
)

type GetAllConsumerResponse struct {
	MetaData MetaData   `json:"metadata"`
	Data     []Consumer `json:"data"`
}

type GetAllBillingResponse struct {
	MetaData MetaData  `json:"metadata"`
	Data     []Billing `json:"data"`
}

type MetaData struct {
	TotalData uint   `json:"total_data"`
	TotalPage uint   `json:"total_page"`
	Page      uint   `json:"page"`
	Limit     uint   `json:"limit"`
	Sort      string `json:"sort"`
	Order     string `json:"order"`
}

type Consumer struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" validate:"required"`
	Nik          string `json:"nik" validate:"required"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	BirthDate    string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	PlaceOfBirth string `json:"place_of_birth" validate:"required"`
	Salary       string `json:"min_salary" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	DtmCrt       time.Time
	DtmUpd       time.Time
}

type ConsumerCreditCard struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Nik            string `json:"nik" validate:"required"`
	ConsumerID     string `json:"consumer_id"`
	CardNumber     string `json:"card_number" `
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
	RequestLimit   string `json:"request_limit"`
	CreditLimit    string `json:"credit_limit"`
	CurrentBalance string `json:"current_balance"`
	DtmCrt         time.Time
	DtmUpd         time.Time
}

type TransactionCreditCard struct {
	ID             int64  `json:"id"`
	ContractNumber string `json:"contract_number"`
	PurchaseAmount string `json:"purchase_amount"`
	ConsumerID     string `json:"consumer_id"`
	TransactionID  string `json:"transaction_id"`
	OTR            string `json:"otr"`
	CVV            string `json:"cvv"`
	AdminFee       string `json:"admin_fee"`
	RequestLimit   string `json:"request_limit"`
	CreditLimit    string `json:"credit_limit"`
	CurrentBalance string `json:"current_balance"`
	DtmCrt         time.Time
	DtmUpd         time.Time

	// Transaction table fields
	TransactionDate  time.Time `json:"transaction_date"`
	InstallmentCount string    `json:"installment_count"`
	InterestAmount   float64   `json:"interest_amount"`
	AssetName        string    `json:"asset_name"`

	// CreditPurchase table fields
	ProductType string `json:"product_type"`
	Status      string `json:"status"`
}
type GetAllCreditCardRequest struct {
	Page           int64  `json:"page"`
	Limit          int64  `json:"limit"`
	Sort           string `json:"sort"`
	Order          string `json:"order"`
	CardNumber     string `json:"card_number" `
	ExpirationDate string `json:"expiration_date"`
	CreditLimit    string `json:"credit_limit"`
	CurrentBalance string `json:"current_balance"`
}
type GetAllCreditCardResponse struct {
	MetaData MetaData             `json:"metadata"`
	Data     []ConsumerCreditCard `json:"data"`
}

type GetAllConsumerRequest struct {
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Email  string `json:"email"`
}

type GetAllBillingRequest struct {
	Page       int64  `json:"page"`
	Limit      int64  `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	BillAmount string `json:"bill_amount"`
	DueDate    string `json:"due_date"`
	Status     string `json:"status"`
}

type Billing struct {
	ID             int       `json:"id"`
	ConsumerID     int       `json:"consumer_id"`
	TransactionID  int       `json:"transaction_id"`
	ContractNumber string    `json:"contract_number"`
	BillAmount     float64   `json:"bill_amount"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `json:"status"`
	DtmCrt         time.Time `json:"dtm_crt"`
	DtmUpd         time.Time `json:"dtm_upd"`
}

type ActionUsecase interface {
	GetConsumer(ctx context.Context, req GetAllConsumerRequest) (res GetAllConsumerResponse, err error)
	GetConsumerDetail(ctx context.Context, req int) (res Consumer, err error)
	PostConsumer(ctx context.Context, req Consumer) (err error)
	PostConsumerCreditCard(ctx context.Context, req ConsumerCreditCard) (err error)
	PostTransactionCredit(ctx context.Context, req TransactionCreditCard, token string) (err error)
	GetCreditCard(ctx context.Context, req GetAllCreditCardRequest) (res GetAllCreditCardResponse, err error)
	GetCreditCardDetail(ctx context.Context, req int) (res ConsumerCreditCard, err error)
	PatchBilling(ctx context.Context, req int, token string) (err error)
	GetBilling(ctx context.Context, req GetAllBillingRequest) (res GetAllBillingResponse, err error)
}

type ActionMySQLRepository interface {
	GetConsumer(ctx context.Context, req GetAllConsumerRequest) (res []Consumer, err error)
	CountConsumer(ctx context.Context, req GetAllConsumerRequest) (res int64, err error)
	CheckValidate(ctx context.Context, coloumn string, value string) (err error)
	CheckValidateCC(ctx context.Context, coloumn string, value string) (count int, err error)
	CheckValidateBilling(ctx context.Context, coloumn string, value string) (count int, err error)
	GetConsumerByParameter(ctx context.Context, coloumn string, value string) (res Consumer, err error)
	PostConsumer(ctx context.Context, req Consumer) (err error)
	GetConsumerDetail(ctx context.Context, req int) (res Consumer, err error)
	PostConsumerCreditCard(ctx context.Context, req ConsumerCreditCard) (err error)
	CountCreditCard(ctx context.Context, req GetAllCreditCardRequest) (res int64, err error)
	GetCreditCard(ctx context.Context, req GetAllCreditCardRequest) (res []ConsumerCreditCard, err error)
	GetCreditCardDetail(ctx context.Context, req int) (res ConsumerCreditCard, err error)
	PostTransactionCredit(ctx context.Context, req TransactionCreditCard) (count int64, err error)
	PostBillingCredit(ctx context.Context, req Billing) (err error)
	UpdateCreditBalance(ctx context.Context, bill float64, id string) (err error)
	CountCreditBalance(ctx context.Context, id string) (res float64, err error)
	PatchBilling(ctx context.Context, id string) (err error)
	GetBilling(ctx context.Context, req GetAllBillingRequest) (res []Billing, err error)
	CountBilling(ctx context.Context, req GetAllBillingRequest) (res int64, err error)
}
