package transaction

import (
	"starup/campaign"
	"starup/user"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID         int
	CampaignID int
	UserId     int
	Amount     int
	Status     string
	Code       string
	User       user.User `gorm:"foreignKey:UserId"` /// UserId ambil dari line atas nya
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
