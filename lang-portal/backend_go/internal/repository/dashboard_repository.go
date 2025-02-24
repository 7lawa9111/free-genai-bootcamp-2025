package repository

import (
    "lang-portal/backend_go/internal/models"
)

type DashboardRepository struct {
    StudySessionRepo *StudySessionRepository
    WordRepo        *WordRepository
    GroupRepo       *GroupRepository
}

func (r *DashboardRepository) GetStats() (*models.DashboardStats, error) {
    // Implementation here
    return &models.DashboardStats{}, nil
} 