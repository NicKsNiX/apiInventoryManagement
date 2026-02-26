# API Inventory Management

A Go-based REST API for inventory management using Fiber framework and MSSQL database.

## Prerequisites

- Go 1.24.4 (for local development)
- Docker & Docker Compose (for containerized deployment)
- MSSQL Server 2019 or later

## Project Structure

```
├── controllers/          # API handlers
├── database/            # Database connection and queries
├── models/              # Data models
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── Dockerfile          # Docker configuration
├── docker-compose.yml  # Docker Compose configuration
└── .env               # Environment configuration (not in repo)
```

## Environment Variables

Create a `.env` file in the root directory with the following configuration:

```
PORT=4002
DB_HOST=your-mssql-host
DB_PORT=1433
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-database-name
```

## Local Development

### Setup

1. Install dependencies:
```bash
go mod download
```

2. Create `.env` file with your database credentials

3. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:4002`

## Docker Deployment

### Build Docker Image

```bash
docker build -t apiinventorymanagement:latest .
```

### Run with Docker

```bash
docker run -p 4002:4002 \
  --env-file .env \
  -v $(pwd)/uploads:/root/uploads \
  apiinventorymanagement:latest
```

### Run with Docker Compose

```bash
docker-compose up -d
```

To stop the containers:
```bash
docker-compose down
```

## API Endpoints

### Authentication
- `POST /api/login` - User login

### Menu & Inventory Info
- `GET /api/getMenuinventory` - Get menu inventory
- `GET /api/getInventoryInfo` - Get inventory information
- `GET /api/inventoryInformation` - Get inventory information
- `GET /api/inventoryInformationReconfirm` - Get inventory info requiring reconfirmation

### Inventory Updates
- `POST /api/updateQtyInventory` - Update quantity
- `POST /api/updateReconfirmQtyInventory` - Update reconfirm quantity
- `POST /api/updateNoconfirmQtyInventory` - Update no-confirm quantity
- `POST /api/insertQtyInventory` - Insert new quantity record
- `POST /api/insertAdjustQtyInventory` - Insert adjustment record
- `POST /api/confirmEditQtyAdjust` - Confirm edit adjustment

### Inventory Checks & History
- `GET /api/getInventoryCheckList` - Get inventory check list
- `GET /api/showHistoryAuditor` - Show auditor history
- `GET /api/getInventoryReconfirmCount` - Get reconfirm count

### Static Files
- `GET /uploads/*` - Serve uploaded files

## Database Schema

The application uses MSSQL with the following models:
- **Account** - User accounts and authentication
- **Inventory** - Inventory items and quantities
- **Menu** - Menu categories and items
- **Permission** - User permissions
- **Position** - User positions/roles

## CI/CD

GitHub Actions workflow is configured in `.github/workflows/docker-build.yml` to automatically build the Docker image on push to `main` branch.

## Image Details

- **Base Image**: Alpine Linux (lightweight)
- **Build Stage**: Go 1.24.4
- **Final Image Size**: ~39.5MB
- **Exposed Port**: 4002

## Troubleshooting

### Database Connection Issues
- Verify MSSQL server is running and accessible
- Check credentials in `.env` file
- Ensure network connectivity between application and database

### Port Already in Use
- Change the PORT in `.env` file
- Or run on different port: `docker run -p 8080:4002 apiinventorymanagement:latest`

### Missing Uploads Directory
- Create the directory: `mkdir uploads`
- Mount volume when running: `-v $(pwd)/uploads:/root/uploads`

## License

This project is proprietary and confidential.

## Repository

GitHub: https://github.com/NicKsNiX/apiInventoryManagement.git
