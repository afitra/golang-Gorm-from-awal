package campaign

import (
	"starup/user"
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CampaignImages   []CampaignImage
	User             user.User
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
