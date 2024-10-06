package command

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"time"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
)

// Envía una respuesta HTTP con un payload en formato JSON.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//Convierte el payload (cualquier estructura de datos) a JSON.
	response, _ := json.Marshal(payload)

	//Establece el tipo de contenido de la respuesta a application/json.
	w.Header().Set("Content-Type", "application/json")
	//Establece el código de estado HTTP de la respuesta.
	w.WriteHeader(code)
	//Escribe el JSON convertido en el cuerpo de la respuesta.
	w.Write(response)
}

// Envía una respuesta de error en formato JSON.
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	//Utiliza RespondWithJSON para enviar el mensaje de error como JSON.
	RespondWithJSON(w, code, map[string]string{"message": msg})
}

// Obtiene la calificación SSL más baja entre los endpoints actuales.
func GetLowestGradeCurrent(data []coreDomain.Endpoint) string {
	//Almacena los valores ASCII de las calificaciones SSL que no sean "A+".
	var gradeASCII []int
	//Almacena la calificación más baja.
	var grade string

	//Itera sobre cada endpoint.
	for _, dataElement := range data {
		//Si la calificación no es "A+", convierte el primer carácter de Grade a su valor ASCII y lo añade a gradeASCII.
		if dataElement.Grade != "A+" {
			gradeASCII = append(gradeASCII, int(dataElement.Grade[0]))
		} else {
			//Si alguna calificación es "A+", asigna directamente "A+" a grade.
			grade = "A+"
		}
	}

	//Si hay calificaciones distintas de "A+", ordena gradeASCII de mayor a menor.
	if len(gradeASCII) > 0 {
		sort.Slice(gradeASCII, func(i, j int) bool {
			return gradeASCII[i] > gradeASCII[j]
		})
		//Convierte el valor ASCII más alto de vuelta a un carácter y lo asigna a grade.
		grade = string(gradeASCII[0])
	}

	return grade
}

// Obtiene la calificación SSL más baja entre los detalles de dominio anteriores.
// La lógica es idéntica a GetLowestGradeCurrent, pero opera sobre un slice de DetailDomain.
func GetLowestGradePrevious(detail []coreDomain.DetailDomain) string {
	var gradeASCII []int
	var grade string

	for _, dataElement := range detail {
		if dataElement.Grade != "A+" {
			gradeASCII = append(gradeASCII, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	if len(gradeASCII) > 0 {
		sort.Slice(gradeASCII, func(i, j int) bool {
			return gradeASCII[i] > gradeASCII[j]
		})

		grade = string(gradeASCII[0])
	}

	return grade
}

// Determina si ha habido un cambio en los servidores o en sus calificaciones SSL.
func ValidateChangeServer(loc *time.Location, payload coreDomain.Domain, data coreDomain.SSL, detailsDomain []coreDomain.DetailDomain, changeServer bool) bool {
	//Calcula la diferencia en horas desde la última consulta.
	hours := DiffHours(loc, payload)
	//Si ha pasado al menos una hora desde la última consulta
	if hours >= 1 {
		//Si el número de endpoints actuales y anteriores es el mismo, compara cada uno.
		if len(data.Endpoints) == len(detailsDomain) {
			for i := 0; i < len(data.Endpoints); i++ {
				//Si hay diferencias en la calificación SSL, nombre del servidor o dirección IP, marca changeServer como true.
				if data.Endpoints[i].Grade != detailsDomain[i].Grade ||
					data.Endpoints[i].ServerName != detailsDomain[i].ServerName ||
					data.Endpoints[i].IpAddress != detailsDomain[i].IpAddress {
					changeServer = true
				}
			}
			//Si el número de endpoints actuales es mayor que el número de detalles anteriores, marca changeServer como true.
		} else if len(data.Endpoints) > len(detailsDomain) {
			changeServer = true
		}
	}
	return changeServer
}

// Calcula la diferencia en horas entre la hora actual y la última consulta.
func DiffHours(loc *time.Location, payload coreDomain.Domain) float64 {
	//Crea una instancia de time.Time para el momento actual.
	t1 := time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		0, 0, 0, loc)

	t2 := time.Date(payload.LastConsult.Year(),
		payload.LastConsult.Month(),
		payload.LastConsult.Day(),
		payload.LastConsult.Hour(),
		0, 0, 0, loc)

	//Crea una instancia de time.Time para el momento de la última consulta.
	return t1.Sub(t2).Hours()
}

// Limpia la dirección URL, eliminando esquemas como http://.
func ValidateURL(address string) string {
	//Compila una expresión regular que coincide con esquemas comunes de URL.
	space := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)`)
	//Elimina los esquemas encontrados en address.
	address = space.ReplaceAllString(address, "")
	return address
}
