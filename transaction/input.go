package transaction

import "starup/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
