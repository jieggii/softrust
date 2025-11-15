# First Step
- The web app: [softthat.fit](https://softthat.fit)
- API documentation (Swagger): [api.softthat.fit/swagger](https://api.softthat.fit/swagger)

## Running

### Development

1. Create file backed.env in the root directory with the following variables:

```env
GENERATIVE_AI_API_KEY=<your key>
```

Run
```bash
docker compose -f docker-compose-dev.yml up --build --attach backend --attach frontend
```