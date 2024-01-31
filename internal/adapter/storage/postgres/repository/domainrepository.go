package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	corePort "github.com/Geovanny0401/clubhub/internal/core/port"
)

func NewSQLDomainRepo(Conn *sql.DB) corePort.DomainRepo {
	return &sqlDomainRepo{
		Conn: Conn,
	}
}

type sqlDomainRepo struct {
	Conn *sql.DB
}

func (sql *sqlDomainRepo) CreateDomain(ctx context.Context, domain coreDomain.Domain) (int64, error) {
	var domainId int64
	query := "INSERT INTO domain(address, last_consult) VALUES($1, $2) RETURNING id"

	stmt, err := sql.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(domain.Address, domain.LastConsult).Scan(&domainId)
	if err != nil {
		return -1, err
	}

	return domainId, err

}

func (sql *sqlDomainRepo) CreateDetailDomain(ctx context.Context, detailDomain coreDomain.DetailDomain) error {
	query := "INSERT INTO detail_domain(id_domain, ipaddress, servername, grade, date) VALUES($1, $2, $3, $4, $5)"

	stmt, err := sql.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Println(err)
	}

	_, err = stmt.ExecContext(ctx, detailDomain.DomainID, detailDomain.IpAddress, detailDomain.ServerName, detailDomain.Grade, detailDomain.Date)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
	}

	return err
}

func (sql *sqlDomainRepo) UpdateLastGetDomain(ctx context.Context, DomainId int64, date time.Time) error {
	query := "UPDATE domain SET last_consult = $1 WHERE id = $2"

	stmt, err := sql.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Println(err)
	}

	_, err = stmt.ExecContext(ctx, date, DomainId)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
	}

	return err
}

func (sql *sqlDomainRepo) GetAllDomain(ctx context.Context) ([]coreDomain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consult FROM domain AS dm"

	rows, err := sql.Conn.QueryContext(ctx, query)

	if err != nil {
		log.Println(err)
	}

	return coreDomain.BuildDomains(rows)
}

func (sql *sqlDomainRepo) GetDomainByAddress(ctx context.Context, address string) (coreDomain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consult FROM domain AS dm WHERE dm.address = $1"

	rows, err := sql.Conn.QueryContext(ctx, query, address)

	if err != nil {
		log.Println(err)
	}

	return coreDomain.BuildDomain(rows)
}

func (sql *sqlDomainRepo) GetDetailsByDomain(ctx context.Context, idDomain int64, countServer int) ([]coreDomain.DetailDomain, error) {
	query := "SELECT dt.id, dt.id_domain, dt.ipaddress, dt.grade, dt.servername, dt.date FROM domain AS dm INNER JOIN detail_domain AS dt ON dt.id_domain = dm.id WHERE dm.id = $1 ORDER BY dt.id, dt.date DESC LIMIT $2"

	rows, err := sql.Conn.QueryContext(ctx, query, idDomain, countServer)

	if err != nil {
		log.Println(err)
	}

	return coreDomain.BuildDetailsDomain(rows)
}
