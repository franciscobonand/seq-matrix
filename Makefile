run:
	go run .

test:
	go test ./... -count=1

coverage:
	go test ./... -count=1 -cover

compose-up:
	docker-compose up

compose-build:
	docker-compose up --build --d

compose-up-mongo:
	docker-compose up mongo mongo-admin

compose-build-mongo:
	docker-compose up --build --d mongo mongo-admin

compose-down:
	docker-compose down

run-loadtest:
	docker run --rm -i --net=host grafana/k6 run - <loadtest/script.js