# Voucher System API

A Go-based REST API for managing voucher redemption system with brands, vouchers, customers, and transactions.

## Features

- **Brand Management**: Create and manage brands
- **Voucher Management**: Create vouchers with point costs and validity periods
- **Customer Management**: Manage customers with point balances
- **Transaction System**: Handle voucher redemptions with multiple vouchers
- **Database Migration**: Automated schema management
- **Unit Testing**: Comprehensive test coverage
- **Validation**: Input validation and error handling

## Database Schema

### Relationships
- One brand can have many vouchers
- One voucher belongs to exactly one brand
- Every voucher has a "cost in point" field
- Customers can make redemption transactions
- Transactions can contain multiple voucher items

### Tables
- `brands`: Brand information
- `vouchers`: Voucher details with point costs
- `customers`: Customer information with point balances
- `transactions`: Redemption transaction records
- `transaction_items`: Individual voucher items in transactions

## API Endpoints

### Brands
- `POST /api/v1/brand` - Create a new brand
- `GET /api/v1/brand` - Get all brands (with pagination)
- `GET /api/v1/brand/:id` - Get a specific brand

### Vouchers
- `POST /api/v1/voucher` - Create a new voucher
- `GET /api/v1/voucher?id={voucher_id}` - Get a specific voucher
- `GET /api/v1/voucher/brand?id={brand_id}` - Get all vouchers by brand
- `GET /api/v1/voucher/all` - Get all vouchers (with pagination)

### Customers
- `POST /api/v1/customer` - Create a new customer
- `GET /api/v1/customer` - Get all customers (with pagination)
- `GET /api/v1/customer/:id` - Get a specific customer
- `PUT /api/v1/customer/:id/points` - Update customer points

### Transactions
- `POST /api/v1/transaction/redemption` - Create a redemption transaction
- `GET /api/v1/transaction/redemption?transactionId={transactionId}` - Get transaction details
- `GET /api/v1/transaction/customer?customerId={customerId}` - Get customer transactions

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd my-backend-app
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp config.env.example config.env
   # Edit config.env with your database credentials
   ```

4. **Create database**
   ```sql
   CREATE DATABASE voucher_system;
   ```

5. **Run database migration**
   ```bash
   # The application will auto-migrate on startup
   # Or manually run the migration file: migrations/001_initial_schema.sql
   ```

## Configuration

Edit `config.env` file:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=voucher_system
SERVER_PORT=8080
```

## Running the Application

### Development Mode
```bash
go run main.go
```

### Production Mode
```bash
GIN_MODE=release go run main.go
```

### Using Docker (if available)
```bash
docker build -t voucher-system .
docker run -p 8080:8080 voucher-system
```

## Testing

### Run all tests
```bash
go test ./...
```

### Run specific test suite
```bash
go test ./tests -v
```

### Run tests with coverage
```bash
go test ./... -cover
```

## API Examples

### Create a Brand
```bash
curl -X POST http://localhost:8080/api/v1/brand \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Indomaret",
    "description": "Indonesian retail chain",
    "logo_url": "https://example.com/indomaret-logo.png"
  }'
```

### Create a Voucher
```bash
curl -X POST http://localhost:8080/api/v1/voucher \
  -H "Content-Type: application/json" \
  -d '{
    "brand_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Discount Voucher",
    "description": "Get 10% off on groceries",
    "cost_in_point": 50000,
    "valid_from": "2024-01-01T00:00:00Z",
    "valid_to": "2024-12-31T23:59:59Z"
  }'
```

### Create a Customer
```bash
curl -X POST http://localhost:8080/api/v1/customer \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890",
    "points": 100000
  }'
```

### Make a Redemption
```bash
curl -X POST http://localhost:8080/api/v1/transaction/redemption \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "550e8400-e29b-41d4-a716-446655440001",
    "items": [
      {
        "voucher_id": "550e8400-e29b-41d4-a716-446655440002",
        "quantity": 2
      }
    ]
  }'
```

### Get Transaction Details
```bash
curl -X GET "http://localhost:8080/api/v1/transaction/redemption?transactionId=550e8400-e29b-41d4-a716-446655440003"
```

## Project Structure

```
my-backend-app/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── config.env              # Environment configuration
├── models/
│   └── models.go           # Database models
├── database/
│   └── database.go         # Database connection and initialization
├── handlers/
│   ├── brand_handler.go    # Brand-related handlers
│   ├── voucher_handler.go  # Voucher-related handlers
│   ├── customer_handler.go # Customer-related handlers
│   └── transaction_handler.go # Transaction-related handlers
├── routes/
│   └── routes.go           # API route definitions
├── migrations/
│   └── 001_initial_schema.sql # Database migration
├── tests/
│   ├── brand_handler_test.go   # Brand handler tests
│   └── voucher_handler_test.go # Voucher handler tests
└── README.md               # This file
```

## Validation Rules

### Brand
- Name: Required, 2-255 characters
- Description: Optional
- Logo URL: Optional, URL format

### Voucher
- Brand ID: Required, valid UUID
- Name: Required, 2-255 characters
- Cost in Point: Required, greater than 0
- Valid From/To: Optional, valid date range

### Customer
- Name: Required, 2-255 characters
- Email: Required, valid email format, unique
- Phone: Optional
- Points: Optional, non-negative

### Transaction
- Customer ID: Required, valid UUID
- Items: Required, non-empty array
- Each item must have valid voucher ID and quantity > 0

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests to ensure everything works
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue in the repository. 