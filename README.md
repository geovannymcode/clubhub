# Clubhub
Prueba técnica para puesto senior de backend en clubhub

## Descripción
Se desarrollo un servicio web RESTful escrito en el lenguaje de programación Go. Este proyecto es parte de una prueba técnica para una vacante de senior de backend en clubhub, donde se necesitaba diseñar y construir un servicio que permita la administración de franquicias hoteleras, debe permitir almacenar, consultar y actualizar la información. Esta prueba contribuyo a mi proceso de aprendizaje en el lenguaje y la comprensión de la Arquitectura Hexagonal implementada en Go.

Utilice [Chi](https://github.com/go-chi/chi) como enrutador HTTP, [PostgreSQL](https://www.postgresql.org/) como base de datos con [pq](https://github.com/lib/pq), Tambien [swagger](https://github.com/swaggo/http-swagger), [swag](https://github.com/swaggo/swag) para generar la documentacion en swagger y [goquery](github.com/PuerkitoBio/goquery) como consultaen documento HTML.

Por otra parte se usaron los siguientes patrones de diseño.
  - **Pattern Adapter**
  - **Dependency Injection**
  - **Pattern Repository**
  - **Pattern Factory**

## Ejecutar Proyecto
En primer lugar se debe clonar el repositorio donde se encuentra el proyecto. Luego ubicarse dentro de la carpeta de clubhub. Después vamos a construir ejecutar Docker. Asegúrese de que tiene docker instalado en su máquina. 


1) Clonamos el repositorio
```bash
git clone https://github.com/Geovanny0401/clubhub
```

2) Nos ubicamos dentro de la carpeta
```bash
 cd clubhub
```

3) Comando para subir Imagen
```bash
   docker-compose up --build -d
```

4) Comando para bajar la Imagen
```bash
    docker-compose down
```

5) Esta paso es opcional, despues de ejecutarlo hacemos el paso # 3 
```bash
chmod +x run.sh
./run.sh
```
## Documentación

La documentación de la API se encuentra en el directorio `docs/`. Para ver la documentación, abra el navegador y vaya a `http://localhost:8001/swagger/index.html`. La documentación se genera utilizando [swaggo](https://github.com/swaggo/swag/) con el middleware [gin-swagger](https://github.com/swaggo/gin-swagger/).

## Diseño
Este el diseño de la Arquitectura Hexagonal en la aplicación como esta compuesta.

![Captura de pantalla 2024-01-31 a las 23 53 52](https://github.com/Geovanny0401/clubhub/assets/10421376/8da6dca2-116d-4918-aec4-4512c69420f0)

## Endpoints

### Todo Address
- Path : `/clubhub`
- Method: `GET`
- Response: `200`
Url: 
```bash
curl --request GET \
  --url 'http://localhost:8001/clubhub?=' \
  --header 'User-Agent: insomnia/8.6.0'
```
  
### Detalle por Address
- Path : `/clubhub/address={address}`
- Method: `GET`
- Response: `200`

Url: 
```bash
curl --request GET \
  --url 'http://localhost:8001/clubhub/address=marriott.com?=' \
  --header 'User-Agent: insomnia/8.6.0' 
```


