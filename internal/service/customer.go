package service

import (
	"context"

	"uuid"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type customerService struct {
	repository repository.CustomerRepository
}

func NewCustomerService(repository repository.CustomerRepository) CustomerService {
	return &customerService{repository: repository}
}

func (s *customerService) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	if customer.ID == "" {
		customer.ID = uuid.NewV4().String()
	}
	return s.repository.CreateCustomer(ctx, customer)
}

func (s *customerService) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	return s.repository.ListCustomers(ctx)
}

func (s *customerService) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	return s.repository.GetCustomer(ctx, id)
}

func (s *customerService) UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error {
	return s.repository.UpdateCustomer(ctx, id, customer)
}

func (s *customerService) DeleteCustomer(ctx context.Context, id string) error {
	return s.repository.DeleteCustomer(ctx, id)
}

var _ CustomerService = (*customerService)(nil)
