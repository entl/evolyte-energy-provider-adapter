# Evolyte Energy Provider Adapter

![Go](https://img.shields.io/badge/Go-1.23-blue?logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo_Framework-Web-blue?logo=go)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Container-2496ED?logo=docker)
![slog](https://img.shields.io/badge/Logging-slog-lightgrey)
![Validator](https://img.shields.io/badge/Validation-go--playground%2Fvalidator-green)

---

A lightweight Go-based web service that acts as an adapter for energy provider APIs, such as Enode. It facilitates secure authentication, token management, and future device integration in Evolyte’s energy management ecosystem.

---

## ✨ Key Features

| Feature                    | Description |
|----------------------------|-------------|
| **Enode API Integration**  | Connects with Enode to authenticate and retrieve access tokens using `client_credentials` flow. |
| **RESTful Interface**      | Clean and extensible HTTP endpoints using the Echo framework. |
| **Graceful Shutdown**      | Manages server shutdown with proper signal handling. |
| **Structured Logging**     | Uses `slog` for consistent, JSON-formatted logs. |
| **Container-Ready**        | Fully dockerized for portable deployment. |

---

## 🧱 Tech Stack

| Component     | Tooling                     |
|---------------|-----------------------------|
| Language      | Go 1.23                     |
| Framework     | Echo v4                     |
| Caching       | Redis                       |
| Validation    | `go-playground/validator`   |
| Configuration | `.env` with `godotenv`      |
| Logging       | Go `slog` (JSON output)     |
| Container     | Docker + Docker Compose     |

---

## ⚙️ Environment Configuration

Create a `.env` file at the root of the project:

```env
PORT=8002
ENODE_CLIENT_ID=your_enode_client_id
ENODE_CLIENT_SECRET=your_enode_client_secret
ENODE_OAUTH_URL=https://oauth.enode.io
ENODE_API_URL=https://api.enode.io
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 🐳 Docker Run

Build and run the service in a container:

```bash
docker compose -f docker-compose.dev.yml --env-file .env.docker up --build
```

---

## 📈 Observability

- Logs are written in structured JSON format via `slog`, captured by stdout (ideal for Filebeat).
- Logs are automatically harvested by the `filebeat` service in the Docker Compose setup based on the `docker-elk` repository.
- Health check is available at `GET /health`.

---


## 📄 License

MIT © 2025 Evolyte Team