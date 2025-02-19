build:
	docker compose -f compose.dev.yaml build

up:
	docker compose -f compose.dev.yaml up

exec:
	docker exec -it app /bin/sh

prod:
	docker compose -f compose.yaml build && docker tag invest-go-app:latest kuuuuya/invest-go:latest && docker push kuuuuya/invest-go:latest && sh deploy.sh