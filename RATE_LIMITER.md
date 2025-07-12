## Como usar

### Com Docker Compose

```bash
docker-compose up --build
```

### Teste manual

```bash
# Teste por IP (limite: 10 req/s)
curl http://localhost:8080/

# Teste com token (limite: 100 req/s)
curl -H "API_KEY: meu-token" http://localhost:8080/
```

### Teste com Apache Bench

```bash
# Teste que deve ativar rate limit
ab -n 15 -c 15 http://localhost:8080/
```