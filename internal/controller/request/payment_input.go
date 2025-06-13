package request

type PaymentInput struct {
	BillingID    int32  `json:"billing_id" example:"2"`
	Amount       int32  `json:"amount" example:"50"`
}
