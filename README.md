# AdRouter

AdRouter is a RESTful API built with Go (Gin framework) that allows users to create, manage, and deliver ad campaigns. It supports features such as CRUD operations on campaigns, targeting management, and ad delivery with precise targeting based on application, country, and operating system. The project is designed for scalability and ease of deployment, leveraging modern DevOps practices and cloud-native technologies.

---

## Features

- **Campaign Management**: Create, update, toggle status, and delete ad campaigns.
- **Targeting Management**: Add and update targeting rules for apps, countries, and operating systems.
- **Delivery**: Retrieve campaigns based on specific targeting criteria with Redis caching for faster responses.
- **Monitoring**: Integrated with Prometheus and Grafana for real-time monitoring and insights.
- **CI/CD**: Automated testing, build, and deployment using GitHub Actions.
- **Scalability**: Deployed as Docker containers on a Kubernetes cluster in Google Cloud Platform (GCP).

---

## Technologies Used

- **Programming Language**: Go (Gin Framework)
- **Database**: PostgreSQL with SQL migrations and sqlc for query generation
- **Caching**: Redis for optimizing delivery route performance
- **Containerization**: Docker
- **Orchestration**: Kubernetes (GCP)
- **Monitoring**: Prometheus and Grafana
- **CI/CD**: GitHub Actions

---

## Folder Structure

```
.
├── .github
│   └── workflows
│       ├── deploy.yml
│       └── test.yml
├── api
│   ├── routes.go
│   └── server.go
├── db
│   ├── migration
│   │   ├── 000001_init_schema.down.sql
│   │   └── 000001_init_schema.up.sql
│   ├── query
│   │   ├── campaign_history.sql
│   │   ├── campaign.sql
│   │   ├── target_app.sql
│   │   ├── target_country.sql
│   │   └── target_os.sql
│   └── sqlc
│       ├── campaign_history_test.go
│       ├── campaign_history.sql.go
│       ├── campaign_test.go
│       ├── campaign.sql.go
│       ├── db.go
│       ├── main_test.go
│       ├── models.go
│       ├── querier.go
│       ├── store_test.go
│       ├── store.go
│       ├── target_app_test.go
│       ├── target_app.sql.go
│       ├── target_country_test.go
│       ├── target_country.sql.go
│       ├── target_os_test.go
│       └── target_os.sql.go
├── templates
│   └── index.html
├── util
│   ├── config_test.go
│   ├── config.go
│   ├── random_test.go
│   └── random.go
├── .gitignore
├── app.env
├── docker-compose.yaml
├── Dockerfile
├── Dockerfile.migrate
├── DOCUMENTATION.md
├── go.mod
├── go.sum
├── main_test.go
├── main.go
├── Makefile
├── migrate.sh
├── README.md
├── sqlc.yaml
└── start.sh
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Kubernetes](https://kubernetes.io/)
- [Go](https://go.dev/)
- PostgreSQL Database
- Redis Server

### Setup Instructions

1. Clone the repository:
    
    ```bash
    git clone https://github.com/vivek-344/AdRouter.git
    cd AdRouter
    ```
    
2. Configure the environment variables in `app.env`:
    
    ```bash
    DB_SOURCE=<your-database-source>
	REDIS_SOURCE=<your-redis-source>
	POSTGRES_HOST=<your-database-host>
	POSTGRES_PASSWORD=<your-database-password>
	POSTGRES_USER=<your-database-user>
	POSTGRES_DB=<your-database-name>
	SERVER_ADDRESS=<your-server-address>
    ```
    
3.  Start the application:
    
    ```bash
    docker compose up
    ```


---

## API Documentation

Comprehensive API documentation is available in [DOCUMENTATION.md](https://github.com/vivek-344/AdRouter/blob/main/DOCUMENTATION.md). It includes details on all endpoints, request parameters, and responses.

---

## CI/CD Pipeline

This project uses GitHub Actions for automated testing, building, and deploying:

1. **Testing**: Executes unit and integration tests.
2. **Building**: Builds Docker images for the application.
3. **Deployment**: Deploys the application to a Kubernetes cluster on GCP.

The workflow files can be found in the `.github/workflows` directory:

- `test.yml`: Handles automated testing.
- `deploy.yml`: Manages building and deployment.

---

## Monitoring

- Prometheus and Grafana are used for monitoring.
- Set up Grafana dashboards to track API usage, errors, and system performance.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b ft/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin ft/your-feature`).
5. Open a pull request.

---

## Contact

For questions or support, please contact the team at [vivekuit344@gmail.com](mailto:vivekuit344@gmail.com).
