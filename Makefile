build:
	docker build -t liveauctioneers/rabbit-amazon-forwarder -f Dockerfile .

push: test build
	docker push liveauctioneers/rabbit-amazon-forwarder

test:
	docker-compose run --rm tests

dev:
	go build
