## Docker

### run command

```sh
docker run -it --rm \
  --name quiz \
  -p 8156:8156 \
  -v ./config/:/app/config/ \
  ghcr.io/flohoss/quiz:latest
```

### compose file

```yml
services:
  quiz:
    image: ghcr.io/flohoss/quiz:latest
    restart: always
    container_name: quiz
    volumes:
      - ./config/:/app/config/
    ports:
      - '8156:8156'
```

## Update Dependencies

```bash
# Node packages
docker compose run --rm --pull always yarn upgrade --latest

# Go packages
docker compose run --rm --pull always go get -u && go mod tidy
```
