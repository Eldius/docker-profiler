
profile:
	go run ./cmd/cli -profile -container dummy-container

plot:
	go run ./cmd/cli -plot -container dummy-container

start-container:
	docker run \
		-m 512m \
		--cpus 2 \
		-d \
		--rm \
		--name dummy-container \
		containerstack/alpine-stress \
		stress \
		--cpu 2 \
		--timeout 60s

stop-container:
	docker kill alpine001
