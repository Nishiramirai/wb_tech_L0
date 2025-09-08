.PHONY: build run test logs stop clean

build:
	@echo "Building the application image..."
	docker compose build

run:
	@echo "Starting the application stack..."
	docker compose up --build -d

send:
	@echo "Running the Go test producer..."
	docker compose run --rm go-test-producer
	
logs:
	@echo "Showing application logs..."
	docker compose logs -f app

stop:
	@echo "Stopping the application stack..."
	docker compose down -v

clean:
	@echo "Cleaning up dangling images and stopped containers..."
	docker image prune -f
	docker container prune -f