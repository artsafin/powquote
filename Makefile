SERVER_IMAGE=powquote-server:local
CLIENT_IMAGE=powquote-client:local

all: build

build: build-server build-client

make-network:
	@docker network create powquote-network 2>/dev/null || true

build-server: make-network
	@docker build --quiet --build-arg "TYPE=server" -t $(SERVER_IMAGE) .

build-client: make-network
	@docker build --quiet --build-arg "TYPE=client" -t $(CLIENT_IMAGE) . 1>/dev/null

run-server: build-server
	@docker stop quoteserver 1>/dev/null || true
	@docker run --rm -d \
		--network="powquote-network" \
		--name="quoteserver" \
		-e "LISTEN=:9999" \
		-e "COMPLEXITY=6" \
		$(SERVER_IMAGE) 1>/dev/null
	@docker logs -f quoteserver

run-client: build-client
	@docker run --rm \
		--network="powquote-network" \
		-e "SERVER=quoteserver:9999" \
		-e "VERBOSE=0" \
		$(CLIENT_IMAGE)