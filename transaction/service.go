package transaction

import (
	"starup/campaign"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {

	return &service{repository, campaignRepository}
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
