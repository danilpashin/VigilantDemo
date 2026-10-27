package service

import (
	"context"

	"uuid"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type cardService struct {
	repository repository.CardRepository
}

func NewCardService(repository repository.CardRepository) CardService {
	return &cardService{repository: repository}
}

func (s *cardService) CreateCard(ctx context.Context, card *models.Card) error {
	if card.ID == "" {
		card.ID = uuid.NewV4().String()
	}
	return s.repository.CreateCard(ctx, card)
}

func (s *cardService) ListCards(ctx context.Context) ([]models.Card, error) {
	return s.repository.ListCards(ctx)
}

func (s *cardService) GetCard(ctx context.Context, id string) (*models.Card, error) {
	return s.repository.GetCard(ctx, id)
}

func (s *cardService) UpdateCard(ctx context.Context, id string, card *models.Card) error {
	return s.repository.UpdateCard(ctx, id, card)
}

func (s *cardService) DeleteCard(ctx context.Context, id string) error {
	return s.repository.DeleteCard(ctx, id)
}

var _ CardService = (*cardService)(nil)
