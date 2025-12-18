.PHONY: help build test docker-build docker-run docker-stop docker-logs docker-push clean migrate up down logs colima-start colima-stop colima-status

# Docker image configuration
IMAGE_NAME = nrzaman/baos-birthday-bot
VERSION ?= latest

help:
	@echo "Baos Birthday Bot - Available Commands:"
	@echo ""
	@echo "Build & Test:"
	@echo "  make build         - Build the Go binary"
	@echo "  make test          - Run all tests"
	@echo "  make migrate       - Run database migration"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build the Docker image (VERSION=v1.0.0 for specific version)"
	@echo "  make docker-run    - Run the bot in Docker"
	@echo "  make docker-stop   - Stop the Docker container"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make docker-push   - Push image to Docker Hub (pushes VERSION and latest)"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make up            - Start with docker-compose"
	@echo "  make down          - Stop docker-compose"
	@echo "  make logs          - View docker-compose logs"
	@echo ""
	@echo "Colima (Docker Runtime):"
	@echo "  make colima-start  - Start Colima"
	@echo "  make colima-stop   - Stop Colima"
	@echo "  make colima-status - Check Colima status"
	@echo ""

build:
	CGO_ENABLED=1 go build -o bot .

build-migrate:
	CGO_ENABLED=1 go build -o migrate ./cmd/migrate

test:
	go test -v ./...

migrate: build-migrate
	@if [ ! -f ./config/birthdays.json ]; then \
		echo "Error: ./config/birthdays.json not found"; \
		exit 1; \
	fi
	@mkdir -p data
	./migrate -json ./config/birthdays.json -db ./data/birthdays.db

docker-build: build
	docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .

docker-run: docker-build migrate
	docker run -d \
		--name birthday-bot \
		--restart unless-stopped \
		-v $(PWD)/data:/app/data \
		-e DISCORD_BIRTHDAY_BOT_TOKEN="$(DISCORD_BIRTHDAY_BOT_TOKEN)" \
		-e DISCORD_BIRTHDAY_CHANNEL_ID="$(DISCORD_BIRTHDAY_CHANNEL_ID)" \
		-e DATABASE_PATH=/app/data/birthdays.db \
		$(IMAGE_NAME):latest

docker-stop:
	docker stop birthday-bot || true
	docker rm birthday-bot || true

docker-logs:
	docker logs -f birthday-bot

docker-push:
	docker push $(IMAGE_NAME):$(VERSION)
	@if [ "$(VERSION)" != "latest" ]; then \
		docker push $(IMAGE_NAME):latest; \
	fi

up: build migrate
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	rm -f bot migrate
	go clean
	docker-compose down || true
	docker stop birthday-bot || true
	docker rm birthday-bot || true
	docker rmi $(IMAGE_NAME):latest || true
	docker rmi $(IMAGE_NAME):$(VERSION) || true

# Colima commands
colima-start:
	colima start

colima-stop:
	colima stop

colima-status:
	colima status