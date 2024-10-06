package command

import (
	"log"
	"net"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

/*
Realiza una consulta WHOIS para una dirección IP y devuelve los resultados en un mapa de cadenas.
La función primero verifica si la dirección IP es válida, realiza la consulta WHOIS utilizando la librería whois,
analiza los resultados con whois-parser, y finalmente estructura la información relevante en un mapa para ser devuelto.
Este proceso implica interactuar con servicios externos para obtener información sobre la dirección IP
y organizar esta información de manera útil para otras partes de la aplicación.
*/
func RunWhoIs(ipAddr string) map[string][]string {
	//net.ParseIP(ipAddr): Intenta convertir la cadena ipAddr en un objeto IP. Si ipAddr no es una dirección IP válida, ParseIP devuelve nil.
	//if ipObj == nil: Comprueba si ipObj es nil, lo que indica una dirección IP inválida. Si es así, registra un error y retorna un mapa vacío.
	//Esto evita que el resto de la función se ejecute con una entrada inválida.
	ipObj := net.ParseIP(ipAddr)
	if ipObj == nil {
		log.Println("Invalid IP Address!")
		return make(map[string][]string)
	}

	//whois.Whois(ipAddr): Realiza una consulta WHOIS para la dirección IP proporcionada.
	//Esta operación puede fallar por varias razones (e.g., problemas de red, dirección IP sin datos WHOIS disponibles),
	//en cuyo caso err no será nil.
	//if err != nil: Si ocurre un error durante la consulta WHOIS, se registra el error y se retorna nil, deteniendo la ejecución de la función.
	result, err := whois.Whois(ipAddr)
	if err != nil {
		log.Println("Error fetching WHOIS data:", err)
		return nil
	}

	//whoisparser.Parse(result): Intenta analizar la cadena de resultado WHOIS cruda en una estructura estructurada (parsedResult).
	//Este proceso puede fallar si el formato del resultado WHOIS no es reconocido por el analizador.
	//if err != nil: Si hay un error durante el análisis, se registra y la función termina retornando nil.
	parsedResult, err := whoisparser.Parse(result)
	if err != nil {
		log.Println("Error parsing WHOIS data:", err)
		return nil
	}
	//make(map[string][]string): Inicializa outPut, un mapa para almacenar los resultados del análisis.
	//Las claves del mapa son cadenas que representan tipos de información (e.g., "DomainName"), y los valores son listas de cadenas.
	var outPut map[string][]string
	outPut = make(map[string][]string)

	//if parsedResult.Domain != nil: Verifica si hay información del dominio disponible en el resultado parseado. Si es así,
	//extrae el nombre del dominio, la fecha de creación y la fecha de expiración, y los almacena en outPut.
	if parsedResult.Domain != nil {
		outPut["DomainName"] = []string{parsedResult.Domain.Name}
		outPut["CreationDate"] = []string{parsedResult.Domain.CreatedDate}
		outPut["ExpirationDate"] = []string{parsedResult.Domain.ExpirationDate}
	}
	//if parsedResult.Registrar != nil: Similar al paso anterior, pero para información del registrador.
	//Si esta información está disponible, extrae el nombre y el correo electrónico del registrador y los almacena en outPut.
	if parsedResult.Registrar != nil {
		outPut["RegistrarName"] = []string{parsedResult.Registrar.Name}
		outPut["RegistrarEmail"] = []string{parsedResult.Registrar.Email}
	}

	//retorna el mapa outPut lleno con la información extraída del resultado WHOIS parseado.
	return outPut
}
