package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTrasactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {

	formatter := CampaignTransactionFormatter{}

	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt

	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {

	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionsFormatter []CampaignTransactionFormatter

	for _, transaction := range transactions {

		formatter := FormatCampaignTransaction(transaction)

		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter

}

func FormatUserTransaction(transaction Transaction) UserTrasactionFormatter {

	formatter := UserTrasactionFormatter{}

	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}

	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageURL = ""

	if len(transaction.Campaign.CampaignImages) > 0 {

		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName

	}

	formatter.Campaign = campaignFormatter

	return formatter

}

func FormatUserTransactions(transactions []Transaction) []UserTrasactionFormatter {

	if len(transactions) == 0 {
		return []UserTrasactionFormatter{}
	}

	var userTransactionsFormatter []UserTrasactionFormatter

	for _, transaction := range transactions {

		formatter := FormatUserTransaction(transaction)

		userTransactionsFormatter = append(userTransactionsFormatter, formatter)
	}

	return userTransactionsFormatter

}
