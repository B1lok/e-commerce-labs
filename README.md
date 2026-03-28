# Weather Forecast Subscription Service

This repository contains a fully implemented **Weather Forecast API service**, which allows users to **subscribe to regular weather updates** for a selected city.

The service follows a predefined API contract (Swagger 2.0) and includes full backend functionality, a subscription HTML page, unit and integration tests, Docker-based deployments, and a complete CI/CD pipeline to Google Cloud Platform (GCP).

---

## Features

- **REST API** that fully complies with the provided Swagger specification.
- **Subscription system** with email confirmation and unsubscribe endpoints.
- **Support for hourly and daily weather updates**.
- **Persisted data in a relational database**.
- **Database migrations** are automatically applied on service startup.
- **Unit tests** with mocks for core services and repositories.
- **Integration tests** using [Testcontainers](https://www.testcontainers.org/).
- **CI/CD Pipeline**:
  - Runs all tests (unit + integration).
  - Builds and pushes a Docker image to **Docker Hub**.
  - Deploys the app on a **GCP Virtual Machine**.
- **Production configuration**:
  - Uses a **secure GCP-hosted database**.
  - Configured **NGINX reverse proxy**.
  - HTTPS enabled via SSL key.
  - **Custom domain name** configured using https://www.namecheap.com/
- **HTML subscription page** for easier end-user access.
  
---

## API Specification

The API follows this contract:

- `GET /weather?city=CityName` - Get current weather for the given city.
- `POST /subscribe` - Subscribe with email, city, and update frequency (hourly or daily).
- `GET /confirm/{token}` - Confirm a new subscription via email token.
- `GET /unsubscribe/{token}` - Unsubscribe via email token.

---

## Testing

- **Unit Tests**: Implemented for core services and repository layers using mocks.
- **Integration Tests**: Use **Testcontainers** to spin up real PostgreSQL and simulate environment.
- **Test Coverage**: Includes subscription flows, weather retrieval, and token validation logic.

---

## CI/CD Pipeline

CI/CD is configured to:

1. Build and test the application using GitHub Actions.
2. Run integration tests with Docker + Testcontainers.
3. Build and push the application image to **Docker Hub**.
4. SSH into a **GCP VM**, pull the new image, and restart the application.

Secrets (API keys, SSH keys, etc.) are stored securely in GitHub Secrets.

---

## Deployment Architecture

| Component       | Technology Used        |
|----------------|------------------------|
| API Backend     | Go                     |
| DB (Prod)       | GCP Cloud SQL (PostgreSQL) |
| HTML Frontend   | Static HTML + CSS      |
| Reverse Proxy   | NGINX                  |
| SSL             | Namecheap Ssl key |
| Domain          | Custom Domain + DNS    |
| Containerization| Docker + Docker Compose |
| CI/CD           | GitHub Actions         |
| Hosting         | Google Cloud Platform  |

---

##  Production Environment Setup

- Custom domain: https://b1lokcode.me/  - hosted application available here
- SSL certificate installed.
- NGINX configured as a reverse proxy.
- Environment variables stored securely on server.
- Database connection is encrypted and secured by GCP firewall rules.

---

## Local Development

To run the Weather Subscription API locally, follow the steps below.

---

### Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Git](https://git-scm.com/)
- Create your api key at https://www.weatherapi.com/
- Create your email key from your google account
---

### Step 1: Clone the Repository

```bash
git clone https://github.com/B1lok/weather-api.git
cd weather-api
```

### Step 2: Create `.env` File

Create a `.env` file in the project root and populate it with the following environment variables:

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=weather
SERVER_HOST=http://localhost:8080/
SERVER_PORT=8080
WEATHER_API_KEY=your-api-key
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email-adress
EMAIL_PASSWORD=your-api-key
```
### Step 3: Start with Docker Compose

Build and run the application using Docker Compose:

```bash
docker-compose up --build
```
Application now available at http://localhost:8080/ with a custom subscription page

### Running Tests with Testcontainers locally

To run integration tests that use **Testcontainers**, make sure the following prerequisites are met:

#### Prerequisites

- Docker is installed and running on your machine.
- Go is installed (version 1.20+ recommended).

#### Run the Tests

Use the following command to run all tests, including those using Testcontainers:

```bash
go test ./... -v
```
