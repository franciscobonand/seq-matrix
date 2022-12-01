test:
	go test ./... -count=1

coverage:
	go test ./... -count=1 -cover

compose-up:
	docker-compose up --build --d

compose-down:
	docker-compose down