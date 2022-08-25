stop:
	docker-compose stop api-devoteam

up:
	docker-compose up -d api-devoteam

up_live:
	docker-compose up api-devoteam

build:
	docker-compose build api-devoteam

tidy:
	docker-compose run --user="root" --rm api-devoteam go mod tidy

verify:
	docker-compose run --rm api-devoteam go mod verify

test:
	docker-compose run --rm api-devoteam bash -c "cd app/; go test ./... -v -p 1"
