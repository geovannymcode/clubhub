package http

import (
	"net/http"
	"time"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	"github.com/go-chi/chi"

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
	address = handlerCommand.ValidateURL(address)

	data, err := GetDataSSl(address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusNoContent, err.Error())
		return
	}

	if "IN_PROGRESS" == data.Status || "DNS" == data.Status || "" == data.Status {
		handlerCommand.RespondWithJSON(w, http.StatusPartialContent, "Try later the server data is not yet available, Thank you!")
		return
	}

	pageTitle, pageLogo, err := GetTitleAndLogo(address)
	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusPartialContent, "Address not found")
		return
	}

	loc, _ := time.LoadLocation("America/Bogota")
	var detailsDomain []coreDomain.DetailDomain
	var changeServer bool

	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)

	if (coreDomain.Domain{}) == payload {
		dm := coreDomain.Domain{}
		dm.Address = address
		dm.LastConsult = time.Now().In(loc)

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)
		if err != nil {
			handlerCommand.RespondWithError(w, http.StatusNoContent, err.Error())
			return
		}

		saveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))
		if err != nil {
			handlerCommand.RespondWithError(w, http.StatusNoContent, err.Error())
			return
		}

		changeServer = handlerCommand.ValidateChangeServer(loc, payload, data, detailsDomain, changeServer)
		if changeServer {
			err = rp.repo.UpdateLastGetDomain(r.Context(), payload.ID, time.Now())
			saveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	if err != nil {
		handlerCommand.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	dataServer := coreUtil.BuildServer(data, detailsDomain, changeServer, pageTitle, pageLogo)
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
		handlerCommand.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	payloadItems := coreUtil.BuilderAddress(payload)
	handlerCommand.RespondWithJSON(w, http.StatusOK, payloadItems)
}

// Method that saves the main details of the domain or server
func saveDetailDomain(data coreDomain.SSL, idDomain int64, rp *Domain, w http.ResponseWriter, r *http.Request) {
	loc, _ := time.LoadLocation("America/Bogota")

	for _, element := range data.Endpoints {
		dt := coreDomain.DetailDomain{}

		dt.DomainID = idDomain
		dt.IpAddress = element.IpAddress
		dt.Grade = element.Grade
		dt.ServerName = element.ServerName
		dt.Date = time.Now().In(loc)

		err := rp.repo.CreateDetailDomain(r.Context(), dt)

		if err != nil {
			handlerCommand.RespondWithError(w, http.StatusNoContent, err.Error())
		}
	}
}
