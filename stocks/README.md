# STOCKS_SERVICE BUILD INFORMATIONS IN DOCKER

## DOCKER INFORMATION
- 

## ENVIRONMENT VARIABLES
- `DB_HOST`: PostgreSQL host - postgres
- `DB_PORT`: PostgreSQL port - 5433
- `DB_USER`: Database username - postgres
- `DB_PASSWORD`: Database password - 12345
- `DB_NAME`: Database name - stocks_service_db
- `HTTP_PORT`: Application port - 8081
- `READ_TIMEOUT`: HTTP read timeout - 15s
- `WRITE_TIMEOUT`: HTTP write timeout - 15s

## API ENDPOINTS
- `POST /stocks/item/add`**Add a new stock item**
- `POST /stocks/item/delete`**Removes stock item**
- `POST /stocks/item/get`**Get stock item by SKU**
- `POST /stocks/list/location`**List stock items by location**