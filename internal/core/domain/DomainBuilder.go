package domain

import "database/sql"

// Construir un objeto Domain a partir de un conjunto de resultados de base de datos (sql.Rows).
func BuildDomain(rows *sql.Rows) (Domain, error) {

	//Declara result como una instancia de Domain que almacenará el resultado final.
	var result Domain
	//Declara resultEmpty, una instancia vacía de Domain para retornar en caso de error.
	resultEmpty := Domain{}
	//Declara err para capturar y manejar errores.
	var err error
	//Itera sobre cada fila en el conjunto de resultados.
	for rows.Next() {
		//Crea una instancia temporal b de Domain para almacenar los datos de la fila actual.
		b := Domain{}
		//Intenta asignar los valores de las columnas de la fila actual a los campos correspondientes en b. Si hay un error, retorna resultEmpty y el error.
		if err := rows.Scan(&b.ID, &b.Address, &b.LastConsult); err != nil {
			return resultEmpty, err
		}
		//Asigna b a result. (Nota: Esto significa que solo se retendrá el último Domain iterado).
		result = b
	}

	//Verifica si hubo algún error durante la iteración de las filas
	if err := rows.Err(); err != nil {
		return resultEmpty, err
	}

	//Retorna el Domain resultante y cualquier error que haya ocurrido.
	return result, err
}

// Construir una lista (slice) de objetos Domain a partir de un conjunto de resultados de base de datos.
func BuildDomains(rows *sql.Rows) ([]Domain, error) {
	//La estructura es similar a BuildDomain, pero acumula todos los dominios en un slice llamado results.
	var results []Domain
	var err error

	for rows.Next() {
		b := Domain{}

		if err := rows.Scan(&b.ID, &b.Address, &b.LastConsult); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, err
}

// Construir una lista de objetos DetailDomain a partir de un conjunto de resultados de base de datos.
func BuildDetailsDomain(rows *sql.Rows) ([]DetailDomain, error) {
	//La estructura es similar a BuildDomains, adaptada para la estructura DetailDomain.
	var results []DetailDomain
	var err error

	for rows.Next() {
		b := DetailDomain{}

		if err := rows.Scan(&b.ID, &b.DomainID, &b.IpAddress, &b.Grade, &b.ServerName, &b.Date); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, err
}
