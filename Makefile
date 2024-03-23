
run:
	go run ./cmd/cli

start-container:
	docker \
		run \
			-d \
			--rm \
			--name alpine001 \
			-m 16m \
			--cpus=0.1 \
			alpine tail -f /dev/null

stop-container:
	docker kill alpine001
