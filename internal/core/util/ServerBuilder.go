package domain

import (
	handlerCommand "github.com/Geovanny0401/clubhub/internal/adapter/handler/command"
	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
)

func BuilderAddress(payload []coreDomain.Domain) coreDomain.Items {
	payloadItems := coreDomain.Items{}
	items := make([]coreDomain.ItemServe, 0)

	for _, element := range payload {
		items = append(items, coreDomain.ItemServe{Address: element.Address})
	}
	payloadItems.Items = items
	return payloadItems
}

func BuildServer(data coreDomain.SSL, detailsDomain []coreDomain.DetailDomain, changeServer bool, pageTitle string, pageLogo string) coreDomain.DataServe {

	currentGrade := handlerCommand.GetLowestGradeCurrent(data.Endpoints)
	var previousGrade string

	if detailsDomain == nil {
		previousGrade = currentGrade
	} else {
		previousGrade = handlerCommand.GetLowestGradePrevious(detailsDomain)
	}

	dataServer := buildData(data, currentGrade, previousGrade, changeServer, pageTitle, pageLogo)

	return dataServer
}

func buildData(data coreDomain.SSL, currentGrade string, previousGrade string, changeServer bool, pageTitle string, pageLogo string) coreDomain.DataServe {
	servers := make([]coreDomain.Serve, 0)
	dataServer := coreDomain.DataServe{}

	for _, dataElement := range data.Endpoints {
		serve := coreDomain.Serve{}
		result := handlerCommand.RunWhoIs(dataElement.IpAddress)

		serve.Address = dataElement.IpAddress
		serve.SslGrade = dataElement.Grade
		serve.Country = result["Country"][0]
		serve.Owner = result["OrgName"][0]

		servers = append(servers, serve)
	}

	dataServer.Serves = servers
	dataServer.ServersChanged = changeServer
	dataServer.SslGrade = currentGrade
	dataServer.PreviousSslGrade = previousGrade
	dataServer.Logo = pageLogo
	dataServer.Title = pageTitle
	dataServer.IsDown = false

	return dataServer
}
