package models

import "time"

type User struct {
	ID   int32   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type Payment struct {
	ID          int32      `gorm:"primaryKey" json:"id"`
	PayerID     int32      `json:"payer_id"`
	BillingID   int32      `json:"-"`
	Amount	    int32     `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type Billing struct {
	ID        	int32      `gorm:"primaryKey" json:"id"`
	PayerID 	int32      `json:"payer_id"`
	ReceiverID 	int32      `json:"receiver_id"`
	Amount    	int32      `json:"amount"`
	CreatedAt 	time.Time  `json:"created_at"`
}
