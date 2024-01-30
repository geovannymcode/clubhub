package http

import (
	"net/http"
	"time"

	"github.com/Geovanny0401/clubhub/internal/core/domain"
	"github.com/go-chi/chi"

	"github.com/Geovanny0401/clubhub/internal/adapter/handler/command"

	"github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres"
	"github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres/repository"
	"github.com/Geovanny0401/clubhub/internal/core/port"
)

func NewServerHandler(db *postgres.DB) *Domain {
	return &Domain{
		repo: repository.NewSQLDomainRepo(db.SQL),
	}
}

type Domain struct {
	repo port.DomainRepo
}

func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	address = command.ValidateURL(address)

	data, err := GetDataSSl(address)
	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, err.Error())
		return
	}

	if "IN_PROGRESS" == data.Status || "DNS" == data.Status || "" == data.Status {
		command.RespondWithJSON(w, http.StatusPartialContent, "Try later the server data is not yet available, Thank you!")
		return
	}

	pageTitle, pageLogo, err := GetTitleAndLogo(address)
	if err != nil {
		command.RespondWithError(w, http.StatusPartialContent, "Address not found")
		return
	}

	loc, _ := time.LoadLocation("America/Bogota")
	var detailsDomain []domain.DetailDomain
	var changeServer bool

	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)

	if (domain.Domain{}) == payload {
		dm := domain.Domain{}
		dm.Address = address
		dm.LastConsult = time.Now().In(loc)

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)
		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
			return
		}

		saveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))
		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
			return
		}

		changeServer = command.ValidateChangeServer(loc, payload, data, detailsDomain, changeServer)
		if changeServer {
			err = rp.repo.UpdateLastGetDomain(r.Context(), payload.ID, time.Now())
			saveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	dataServer := domain.BuildServer(data, detailsDomain, changeServer, pageTitle, pageLogo)
	command.RespondWithJSON(w, http.StatusOK, dataServer)
}

func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	payloadItems := domain.BuilderAddress(payload)
	command.RespondWithJSON(w, http.StatusOK, payloadItems)
}

func saveDetailDomain(data domain.SSL, idDomain int64, rp *Domain, w http.ResponseWriter, r *http.Request) {
	loc, _ := time.LoadLocation("America/Bogota")

	for _, element := range data.Endpoints {
		dt := domain.DetailDomain{}

		dt.DomainID = idDomain
		dt.IpAddress = element.IpAddress
		dt.Grade = element.Grade
		dt.ServerName = element.ServerName
		dt.Date = time.Now().In(loc)

		err := rp.repo.CreateDetailDomain(r.Context(), dt)

		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
		}
	}
}
