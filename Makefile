build:
	docker compose -f compose.dev.yaml build

up:
	docker compose -f compose.dev.yaml up

exec:
	docker exec -it app /bin/bash