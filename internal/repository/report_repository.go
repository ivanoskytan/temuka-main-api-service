package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, report *model.Report) error
	UpdateReport(ctx context.Context, report *model.Report) error
	GetReportsByCommunity(ctx context.Context, communityId int) ([]*model.Report, error)
	DeleteReport(ctx context.Context, id int) error
}

type ReportRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewReportRepository(db database.PostgresWrapper) ReportRepository {
	return &ReportRepositoryImpl{
		db: db,
	}
}

func (r *ReportRepositoryImpl) CreateReport(ctx context.Context, report *model.Report) error {
	if err := r.db.DB.WithContext(ctx).Create(report).Error; err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}
	return nil
}

func (r *ReportRepositoryImpl) UpdateReport(ctx context.Context, report *model.Report) error {
	if err := r.db.DB.WithContext(ctx).Save(report).Error; err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}
	return nil
}

func (r *ReportRepositoryImpl) GetReportsByCommunity(ctx context.Context, communityId int) ([]*model.Report, error) {
	var reports []*model.Report
	if err := r.db.DB.WithContext(ctx).Where("community_id = ?", communityId).Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to get reports by community: %w", err)
	}
	return reports, nil
}

func (r *ReportRepositoryImpl) DeleteReport(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.Report{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}
	return nil
}
