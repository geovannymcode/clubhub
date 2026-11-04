# Clubhub
Prueba técnica para puesto senior de backend en clubhub

## Descripción
Se desarrollo un servicio web RESTful escrito en el lenguaje de programación Go. Este proyecto es parte de una prueba técnica para una vacante de senior de backend en clubhub, donde se necesitaba diseñar y construir un servicio que permita la administración de franquicias hoteleras, debe permitir almacenar, consultar y actualizar la información. Esta prueba contribuyo a mi proceso de aprendizaje en el lenguaje y la comprensión de la Arquitectura Hexagonal implementada en Go.

Utilice [Chi v5](https://github.com/go-chi/chi) como enrutador HTTP, [PostgreSQL](https://www.postgresql.org/) como base de datos con [pq](https://github.com/lib/pq), [swag](https://github.com/swaggo/swag) para generar la documentación OpenAPI y [goquery](github.com/PuerkitoBio/goquery) como consultaen documento HTML.

Por otra parte se usaron los siguientes patrones de diseño.
  - **Pattern Adapter**
  - **Dependency Injection**
  - **Pattern Repository**
  - **Pattern Factory**

## Ejecutar Proyecto

Se requiere Go 1.26.10 o superior. La imagen Docker incluye esta versión con correcciones de seguridad.

1. Copie `.env.example` a `.env` y configure `DB_ADMIN_PASSWORD`, `DB_PASSWORD` y `API_TOKEN` con secretos independientes generados con `openssl rand -hex 32`.
2. Proteja el archivo con `chmod 600 .env`.
3. Para una instalación nueva, ejecute `docker compose up --build -d`.
4. Envíe `Authorization: Bearer <API_TOKEN>` en cada petición.

La API está disponible solo en `127.0.0.1:8001`; PostgreSQL no publica puertos. Antes de actualizar una instalación con datos existentes, siga la migración de credenciales en [SECURITY.md](SECURITY.md). Allí se documentan también los límites, el uso de HTTPS y las comprobaciones de seguridad.

## Documentación

Los archivos OpenAPI están en `docs/` para consulta local. Swagger no se publica por HTTP.

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
  --header "Authorization: Bearer ${API_TOKEN}"
```
  
### Detalle por Address
- Path : `/clubhub/address={address}`
- Method: `GET`
- Response: `200`

Url: 
```bash
curl --request GET \
  --url 'http://localhost:8001/clubhub/address=marriott.com?=' \
  --header "Authorization: Bearer ${API_TOKEN}"
```


