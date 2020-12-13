package transaction

import (
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
	User       user.User `gorm:"foreignKey:id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
