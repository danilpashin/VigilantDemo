package service

import (
	"context"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type auditLogService struct {
	repository repository.AuditLogRepository
}

func NewAuditLogService(repository repository.AuditLogRepository) AuditLogService {
	return &auditLogService{repository: repository}
}

func (s *auditLogService) CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error {
	return s.repository.CreateAuditLog(ctx, auditLog)
}

func (s *auditLogService) ListAuditLogs(ctx context.Context) ([]models.AuditLog, error) {
	return s.repository.ListAuditLogs(ctx)
}

var _ AuditLogService = (*auditLogService)(nil)
