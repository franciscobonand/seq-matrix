run:
	go run .

test:
	go test ./... -count=1

coverage:
	go test ./... -count=1 -cover

compose-up:
	docker-compose up --build --d

compose-up-mongo:
	docker-compose up --build --d mongo mongo-admin

compose-down:
	docker-compose down