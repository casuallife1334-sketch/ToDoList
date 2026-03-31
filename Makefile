include .env
export

postgres-image:
	docker run -e POSTGRES_PASSWORD=Subs22848 -p 5432:5432 -v ./out/pgdata:/var/lib/postgresql -d postgres:18-bookworm

migrate-create:
	migrate create -ext sql -dir migrations -seq init

migrate-up:
	migrate -path migrations -database ${CONN_STRING} up

migrate-down:
	migrate -path migrations -database ${CONN_STRING} down