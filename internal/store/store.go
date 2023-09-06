package store

import (
	coffeeco "coffeeco/internal"
	"context"

	"github.com/google/uuid"
)

type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}

type Service struct {
	storeRepo Repository
}

func NewService(repo Repository) *Service {
	return &Service{storeRepo: repo}
}

func (s Service) GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	dis, err := s.storeRepo.GetStoreDiscount(ctx, storeID)
	if err != nil {
		return 0, err
	}
	return float32(dis), nil
}
