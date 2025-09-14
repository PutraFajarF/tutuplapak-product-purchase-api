# TutuPlapak Product Purchase API

## Environment Configuration

This application uses environment variables for configuration. You can set them in several ways:

### 1. Using Docker Compose environment variables
The `docker-compose.yml` includes default values that can be overridden.

### 2. Using system environment variables
```bash
export JWT_SECRET="your-secure-secret"
export POSTGRESQL_URL="your-database-url"
# ... other variables
```

### 3. Using .env.example as reference
The `.env.example` file contains all available environment variables and their default values for reference.

## JWT Configuration

The JWT validation middleware requires the following:

- **JWT_SECRET**: Secret key for signing/verifying tokens (change from default!)
- **Issuer**: Must be `"tutuplapak-api"`
- **Audience**: Must be `"tutuplapak-services"`
- **Algorithm**: HS256
- **Claims**: `sub` (user ID), `email`, `phone` (custom claims)

### Example JWT Token Structure
```json
{
  "sub": "user-123",
  "email": "user@example.com",
  "phone": "+6281234567890",
  "iss": "tutuplapak-api",
  "aud": "tutuplapak-services",
  "exp": 1640995200,
  "iat": 1640991600
}
```

## Running the Application

### Using Docker Compose
```bash
# Development (uses default values from docker-compose.yml)
docker-compose up

# Production (override environment variables)
JWT_SECRET="your-production-secret" POSTGRESQL_URL="prod-db-url" docker-compose up

# Override multiple variables
JWT_SECRET="secret" POSTGRESQL_URL="db-url" LOG_LEVEL="info" docker-compose up
```

### Using Go directly
```bash
# Install dependencies
go mod tidy

# Run with environment variables
JWT_SECRET="your-secret" POSTGRESQL_URL="your-db-url" go run cmd/main.go

# Or set environment variables first
export JWT_SECRET="your-secret"
export POSTGRESQL_URL="your-db-url"
go run cmd/main.go
```

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/product` - Create product (requires JWT)
- `GET /api/v1/product` - List products
- `PUT /api/v1/product/{productId}` - Update product (requires JWT)
- `DELETE /api/v1/product/{productId}` - Delete product (requires JWT)
- `POST /api/v1/purchase` - Create purchase
- `POST /api/v1/purchase/{purchaseId}` - Upload proof

## Database Migrations

```
-- Drop tables in correct order (because of FK dependencies)
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS categories CASCADE;

-- Create authentications
CREATE TABLE IF NOT EXISTS authentications (
  id VARCHAR PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Create categories
CREATE TABLE IF NOT EXISTS categories (
  type VARCHAR(32) PRIMARY KEY
);

-- Create products
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  auth_id VARCHAR NOT NULL,
  type VARCHAR(32) NOT NULL,
  name VARCHAR(32) NOT NULL CHECK (char_length(name) >= 4 AND char_length(name) <= 32),
  qty INT NOT NULL CHECK (qty >= 1),
  price INT NOT NULL CHECK (price >= 100),
  sku VARCHAR(32) NOT NULL CHECK (char_length(sku) <= 32),
  file_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_auth FOREIGN KEY (auth_id) REFERENCES authentications (id) ON DELETE CASCADE,
  CONSTRAINT fk_category FOREIGN KEY (type) REFERENCES categories (type),
  CONSTRAINT fk_file FOREIGN KEY (file_id) REFERENCES files (id),
  CONSTRAINT unique_auth_sku UNIQUE (auth_id, sku)
);

-- Isi enum category dari Product Category Table
INSERT INTO categories (type) VALUES
  ('Food'),
  ('Beverage'),
  ('Clothes'),
  ('Furniture'),
  ('Tools');
```