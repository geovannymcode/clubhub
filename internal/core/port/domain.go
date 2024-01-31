package port

import (
	"context"
	"time"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
)

type DomainRepo interface {
	CreateDomain(ctx context.Context, domain coreDomain.Domain) (int64, error)
	CreateDetailDomain(ctx context.Context, detailDomain coreDomain.DetailDomain) error
	UpdateLastGetDomain(ctx context.Context, idDomain int64, date time.Time) error
	GetAllDomain(ctx context.Context) ([]coreDomain.Domain, error)
	GetDomainByAddress(ctx context.Context, address string) (coreDomain.Domain, error)
	GetDetailsByDomain(ctx context.Context, idDomain int64, countServer int) ([]coreDomain.DetailDomain, error)
}
