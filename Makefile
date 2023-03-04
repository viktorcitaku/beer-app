fast-run:
	go run cmd/beer/main.go

build:
	docker build -t viktorcitaku/beer-app -f build/deploy/Dockerfile .

run:
	docker run -it -d -p 3000:3000 viktorcitaku/beer-app

compose-up:
	docker compose --file deployments/docker-compose.yml up --build --detach

compose-down:
	docker compose --file deployments/docker-compose.yml down
