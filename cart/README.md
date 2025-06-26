# CART_SERVICE INFORMATIONS

## DOCKER INFORMATION
- **Image Name**```gadamcuma/cart_service```
- **Tag**```hw7_v3```
- **App Port**```8080```

## ENVIRONMENT VARIABLES
- `DB_HOST`: PostgreSQL host - postgres
- `DB_PORT`: PostgreSQL port - 5433
- `DB_USER`: Database username - postgres
- `DB_PASSWORD`: Database password - 12345
- `DB_NAME`: Database name - cart_service_db
- `HTTP_PORT`: Application port - 8080
- `READ_TIMEOUT`: HTTP read timeout - 15s
- `WRITE_TIMEOUT`: HTTP write timeout - 15s
- `STOCKS_SERVICE_URL`: http://stocks_service_backend:8081

## API ENDPOINTS
- `POST /cart/item/add`**Add a new cart item**
- `POST /cart/item/delete`**Removes cart item by sku and user**
- `POST /cart/list`**List carts of user by id**
- `POST /cart/clear`**Removes all cart items for user**