# Go Backend for Sonar

A simple Go backend API built with Echo, GORM, and SQLite. The project exposes CRUD endpoints for products, categories, carts, and orders, and includes SonarQube/SonarCloud configuration for code quality analysis.

## Tech Stack

- Go 1.24.2
- Echo v4
- GORM
- SQLite
- SonarQube/SonarCloud

## Project Structure

```text
.
├── controllers/              # HTTP route handlers
├── database/                 # Database connection setup
├── models/                   # GORM models and query scopes
├── .github/workflows/        # CI workflow for SonarQube scan
├── main.go                   # Application entry point
├── go.mod                    # Go module definition
├── sonar-project.properties  # Sonar project configuration
└── test.db                   # SQLite database file
```

## Features

- Product management
- Category management
- Cart item management
- Order management
- SQLite database with automatic GORM migrations
- Query filtering for products, categories, carts, and orders
- SonarQube/SonarCloud scan workflow on pull requests and pushes to `main`

## Getting Started

### Prerequisites

- Go 1.24.2 or newer

### Install Dependencies

```bash
go mod download
```

### Run the Application

```bash
go run main.go
```

The server starts on:

```text
http://localhost:8080
```

The application uses SQLite and connects to `test.db`. Database tables are automatically migrated on startup.

## API Endpoints

### Products

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/products` | Create a product |
| `GET` | `/products` | List products |
| `GET` | `/products/:id` | Get a product by ID |
| `PUT` | `/products/:id` | Update a product |
| `DELETE` | `/products/:id` | Delete a product |

Product list query parameters:

- `search`
- `min_price`
- `max_price`

### Categories

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/categories` | Create a category |
| `GET` | `/categories` | List categories |
| `GET` | `/categories/:id` | Get a category by ID |
| `PUT` | `/categories/:id` | Update a category |
| `DELETE` | `/categories/:id` | Delete a category |

Category list query parameters:

- `search`
- `min_products`
- `max_products`

### Cart

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/carts` | Add an item to the cart |
| `GET` | `/carts` | List cart items |
| `DELETE` | `/carts/:id` | Delete a cart item |

Cart list query parameters:

- `product_id`
- `min_quantity`
- `max_quantity`

### Orders

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/orders` | Create an order |
| `GET` | `/orders` | List orders |
| `GET` | `/orders/:id` | Get an order by ID |
| `PUT` | `/orders/:id` | Update order status |
| `DELETE` | `/orders/:id` | Delete an order |

Order list query parameters:

- `status`
- `user_id`

## Example Requests

Create a category:

```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics"}'
```

Create a product:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"Name":"Phone","Description":"Smartphone","Price":999.99,"CategoryID":1}'
```

Add an item to the cart:

```bash
curl -X POST http://localhost:8080/carts \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'
```

Create an order:

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "unit_price": 999.99
      }
    ]
  }'
```

## SonarQube

The repository includes `sonar-project.properties` and a GitHub Actions workflow that runs a SonarQube scan on:

- Pushes to `main`
- Pull requests

The workflow requires the `SONAR_TOKEN` secret to be configured in GitHub.
