package service

import (
	"context"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type rulesetService struct {
	repository repository.RulesetRepository
}

func NewRulesetService(repository repository.RulesetRepository) RulesetService {
	return &rulesetService{repository: repository}
}

func (s *rulesetService) CreateRuleset(ctx context.Context, ruleset *models.Ruleset) error {
	return s.repository.CreateRuleset(ctx, ruleset)
}

func (s *rulesetService) ListRulesets(ctx context.Context) ([]models.Ruleset, error) {
	return s.repository.ListRulesets(ctx)
}

var _ RulesetService = (*rulesetService)(nil)
