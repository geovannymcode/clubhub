package domain

import "github.com/Geovanny0401/clubhub/internal/adapter/handler/command"

func BuilderAddress(payload []Domain) Items {
	payloadItems := Items{}
	items := make([]ItemServer, 0)

	for _, element := range payload {
		items = append(items, ItemServer{Address: element.Address})
	}
	payloadItems.Items = items
	return payloadItems
}

func BuildServer(data SSL, detailsDomain []DetailDomain, changeServer bool, pageTitle string, pageLogo string) DataServer {

	currentGrade := command.GetLowestGradeCurrent(data.Endpoints)
	var previousGrade string

	if detailsDomain == nil {
		previousGrade = currentGrade
	} else {
		previousGrade = command.GetLowestGradePrevious(detailsDomain)
	}

	dataServer := buildData(data, currentGrade, previousGrade, changeServer, pageTitle, pageLogo)

	return dataServer
}

func buildData(data SSL, currentGrade string, previousGrade string, changeServer bool, pageTitle string, pageLogo string) DataServer {
	servers := make([]Server, 0)
	dataServer := DataServer{}

	for _, dataElement := range data.Endpoints {
		serve := Server{}
		result := command.RunWhoIs(dataElement.IpAddress)

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
