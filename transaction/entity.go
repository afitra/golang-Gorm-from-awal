package transaction

import (
	"starup/campaign"
	"starup/user"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User /// UserId ambil dari line atas nya
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
