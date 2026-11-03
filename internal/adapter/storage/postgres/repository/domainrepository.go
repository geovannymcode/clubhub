package repository

import (
	"context"
	"database/sql"
	"time"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	corePort "github.com/Geovanny0401/clubhub/internal/core/port"
)

func NewSQLDomainRepo(conn *sql.DB) corePort.DomainRepo { return &sqlDomainRepo{Conn: conn} }

type sqlDomainRepo struct{ Conn *sql.DB }

func (repo *sqlDomainRepo) CreateDomain(ctx context.Context, domain coreDomain.Domain) (int64, error) {
	var id int64
	err := repo.Conn.QueryRowContext(ctx, "INSERT INTO domain(address, last_consult) VALUES($1, $2) RETURNING id", domain.Address, domain.LastConsult).Scan(&id)
	return id, err
}

func (repo *sqlDomainRepo) CreateDetailDomain(ctx context.Context, detail coreDomain.DetailDomain) error {
	_, err := repo.Conn.ExecContext(ctx, "INSERT INTO detail_domain(id_domain, ipaddress, servername, grade, date) VALUES($1, $2, $3, $4, $5)", detail.DomainID, detail.IpAddress, detail.ServerName, detail.Grade, detail.Date)
	return err
}

func (repo *sqlDomainRepo) UpdateLastGetDomain(ctx context.Context, id int64, date time.Time) error {
	_, err := repo.Conn.ExecContext(ctx, "UPDATE domain SET last_consult = $1 WHERE id = $2", date, id)
	return err
}

func (repo *sqlDomainRepo) GetAllDomain(ctx context.Context) ([]coreDomain.Domain, error) {
	rows, err := repo.Conn.QueryContext(ctx, "SELECT id, address, last_consult FROM domain ORDER BY id DESC LIMIT 1000")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return coreDomain.BuildDomains(rows)
}

func (repo *sqlDomainRepo) GetDomainByAddress(ctx context.Context, address string) (coreDomain.Domain, error) {
	rows, err := repo.Conn.QueryContext(ctx, "SELECT id, address, last_consult FROM domain WHERE address = $1 LIMIT 1", address)
	if err != nil {
		return coreDomain.Domain{}, err
	}
	defer rows.Close()
	return coreDomain.BuildDomain(rows)
}

func (repo *sqlDomainRepo) GetDetailsByDomain(ctx context.Context, id int64, count int) ([]coreDomain.DetailDomain, error) {
	rows, err := repo.Conn.QueryContext(ctx, "SELECT id, id_domain, ipaddress, grade, servername, date FROM detail_domain WHERE id_domain = $1 ORDER BY date DESC, id DESC LIMIT $2", id, min(count, 100))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return coreDomain.BuildDetailsDomain(rows)
}
