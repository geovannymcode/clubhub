package http

import (
	"net/http"
	"time"

	"github.com/Geovanny0401/clubhub/internal/security"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	"github.com/go-chi/chi/v5"

	handlerCommand "github.com/Geovanny0401/clubhub/internal/adapter/handler/command"

	storagePostgres "github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres"
	postgresRespository "github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres/repository"
	corePort "github.com/Geovanny0401/clubhub/internal/core/port"
	coreUtil "github.com/Geovanny0401/clubhub/internal/core/util"
)

// NewServerHandler creates a new ServerHandler instance
func NewServerHandler(db *storagePostgres.DB) *Domain {
	return &Domain{
		repo: postgresRespository.NewSQLDomainRepo(db.SQL),
	}
}

// Domain wraps a DomainRepo interface for dependency injection or repository abstraction.
type Domain struct {
	repo corePort.DomainRepo
}

// @summary We Get By Address
// @description GetByAddress We get all the information by address
// @accept	json
// @produce	json
// @Produce json
// @Param   address path string true "address"
// @Router /clubhub/address/{address} [get]
func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	address, err := security.Domain(address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusBadRequest, "Invalid public domain")
		return
	}

	data, err := GetDataSSl(r.Context(), address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusBadGateway, "External service unavailable")
		return
	}

	if "IN_PROGRESS" == data.Status || "DNS" == data.Status || "" == data.Status {
		handlerCommand.RespondWithJSON(w, http.StatusPartialContent, "Try later the server data is not yet available, Thank you!")
		return
	}

	pageTitle, pageLogo, err := GetTitleAndLogo(r.Context(), address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusPartialContent, "Address not found")
		return
	}

	loc := time.FixedZone("America/Bogota", -5*60*60)
	var detailsDomain []coreDomain.DetailDomain
	var changeServer bool

	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
		return
	}

	if (coreDomain.Domain{}) == payload {
		dm := coreDomain.Domain{}
		dm.Address = address
		dm.LastConsult = time.Now().In(loc)

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)
		if err != nil {
			handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
			return
		}

		if err := saveDetailDomain(data, idDomain, rp, r); err != nil {
			handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
			return
		}
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))
		if err != nil {
			handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
			return
		}

		changeServer = handlerCommand.ValidateChangeServer(loc, payload, data, detailsDomain, changeServer)
		if changeServer {
			err = rp.repo.UpdateLastGetDomain(r.Context(), payload.ID, time.Now())
			if err != nil {
				handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
				return
			}
			if err := saveDetailDomain(data, payload.ID, rp, r); err != nil {
				handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
				return
			}
		}
	}

	dataServer := coreUtil.BuildServer(r.Context(), data, detailsDomain, changeServer, pageTitle, pageLogo)
	if r.Context().Err() != nil {
		handlerCommand.RespondWithError(w, http.StatusGatewayTimeout, "Request timed out")
		return
	}
	handlerCommand.RespondWithJSON(w, http.StatusOK, dataServer)
}

// @summary Get All Addresses
// @description GetAllAddress We get all the addresses we have consulted.
// @accept  json
// @produce  json
// @router /clubhub [get]
func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusInternalServerError, "Request could not be completed")
		return
	}

	payloadItems := coreUtil.BuilderAddress(payload)
	handlerCommand.RespondWithJSON(w, http.StatusOK, payloadItems)
}

// Method that saves the main details of the domain or server
func saveDetailDomain(data coreDomain.SSL, idDomain int64, rp *Domain, r *http.Request) error {
	loc := time.FixedZone("America/Bogota", -5*60*60)

	for _, element := range data.Endpoints {
		dt := coreDomain.DetailDomain{}

		dt.DomainID = idDomain
		dt.IpAddress = element.IpAddress
		dt.Grade = element.Grade
		dt.ServerName = element.ServerName
		dt.Date = time.Now().In(loc)

		err := rp.repo.CreateDetailDomain(r.Context(), dt)

		if err != nil {
			return err
		}
	}
	return nil
}
