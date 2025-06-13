package request

type BillingInput struct {
	PayerID    int32 `json:"payer_id" example:"1"`
	ReceiverID int32 `json:"receiver_id" example:"2"`
}