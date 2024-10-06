package domain

import (
	handlerCommand "github.com/Geovanny0401/clubhub/internal/adapter/handler/command"
	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
)

// Crear un objeto Items que contiene un slice de ItemServe basado en un slice de Domain.
func BuilderAddress(payload []coreDomain.Domain) coreDomain.Items {
	//Inicializa payloadItems para almacenar el resultado.
	payloadItems := coreDomain.Items{}
	//Inicializa un slice vacío de ItemServe.
	items := make([]coreDomain.ItemServe, 0)
	//Itera sobre cada Domain en el payload.
	for _, element := range payload {
		//Añade un nuevo ItemServe a items por cada Domain, usando la dirección de cada Domain.
		items = append(items, coreDomain.ItemServe{Address: element.Address})
	}
	//Asigna el slice de ItemServe a payloadItems.
	payloadItems.Items = items
	//Retorna el objeto Items resultante.
	return payloadItems
}

// Construir un objeto DataServe a partir de la información SSL, detalles del dominio, y otros datos.
// La función primero calcula las calificaciones SSL actuales y anteriores, luego construye la información del servidor usando buildData.
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

// Auxiliar de BuildServer para construir el objeto DataServe.
// Esta función crea un slice de Serve basado en los endpoints, asigna varios campos y retorna el DataServe construido.
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
