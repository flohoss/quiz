## Update Dependencies

```bash
# Node packages
docker compose run --rm --pull always yarn upgrade --latest

# Go packages
docker compose run --rm --pull always go get -u && go mod tidy
```
