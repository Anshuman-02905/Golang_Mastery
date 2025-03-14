# Define default Postgres container name and database
DB_CONTAINER=postgres_stock
DB_NAME=stock_db
DB_USER=root

# Start services in detached mode
dc-up:
	docker compose up -d

# Stop and remove services
dc-down:
	docker compose down

# Connect to PostgreSQL database in the container
psql:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

# Run database migrations (assuming you're using a migration tool like migrate)
# migrate-up:
# 	migrate -path migrations -database "postgresql://$(DB_USER):root@localhost:5432/$(DB_NAME)?sslmode=disable" up

# # Rollback last migration
# migrate-down:
# 	migrate -path migrations -database "postgresql://$(DB_USER):root@localhost:5432/$(DB_NAME)?sslmode=disable" down

# Create stock table
db-create:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "CREATE TABLE IF NOT EXISTS stock (stockid BIGSERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, price VARCHAR(255) NOT NULL, company VARCHAR(255) NOT NULL);"

# Drop stock table
db-drop:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP TABLE IF EXISTS stock;"
