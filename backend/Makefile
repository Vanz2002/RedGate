USERNAME = fteen
PASSWORD = fteen
PREFIX = red-gate
DB_NAME = $(PREFIX)-db
DB_NAME_TEST = $(PREFIX)-db-test
IMAGE_NAME = my-postgres:latest
CONTAINER_NAME = $(PREFIX)-container

#
run:
	make pg-start && \
	go run main.go

# FIRST START
pg-run:
	docker run -d --name $(CONTAINER_NAME) -p 5432:5432 $(IMAGE_NAME)
pg-createcred:
	docker exec -it $(CONTAINER_NAME) createuser -U postgres $(USERNAME) && \
	docker exec -it $(CONTAINER_NAME) bash -c "su - postgres -c 'psql -c \"ALTER USER $(USERNAME) WITH PASSWORD '\''$(PASSWORD)'\'';\"'"
pg-createdb:
	docker exec -it $(CONTAINER_NAME) createdb -U postgres -O $(USERNAME) $(DB_NAME)
pg-createdb-test:
	docker exec -it $(CONTAINER_NAME) createdb -U postgres -O $(USERNAME) $(DB_NAME_TEST)

# AFTER 
pg-start:
	docker start $(CONTAINER_NAME)
pg-stop:
	docker stop $(CONTAINER_NAME)
pg-dropuser:
	# make pg-dropdb && \
	docker exec -it $(CONTAINER_NAME) dropuser -U postgres $(USERNAME)
pg-dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb -U postgres $(DB_NAME)
pg-psql:
	docker exec -it $(CONTAINER_NAME) psql -U $(USERNAME) -d $(DB_NAME)
pg-psql-test:
	docker exec -it $(CONTAINER_NAME) psql -U $(USERNAME) -d $(DB_NAME_TEST)


# db migration
migrate-init:
	migrate create -ext sql -dir db/migration -seq db_scheme
migrate-up:
	migrate -path db/migration -database "postgresql://$(USERNAME):$(PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose up
migrate-up-test:
	migrate -path db/migration -database "postgresql://$(USERNAME):$(PASSWORD)@localhost:5432/$(DB_NAME_TEST)?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migration -database "postgresql://$(USERNAME):$(PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose down
migrate-fix:
	migrate -path db/migration -database "postgresql://$(USERNAME):$(PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" force VERSION