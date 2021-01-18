package transaction

import (
	"starup/campaign"
	"starup/payment"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {

	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil

}

func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {

		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {

	transaction := Transaction{}

	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	// transaction.Code = fmt.Sprintf("TS/%d-%s-%s-%s", userID, "campaign", currentTime.Format("2006-01-02-3:4:5"), file.Filename)

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {

		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{

		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {

		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {

		return newTransaction, err
	}

	return newTransaction, nil
}
