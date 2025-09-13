### DB MIgrations

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