# Another World Shop

A simple Go-based backend project for managing products, orders, users, and more. Designed for easy setup and clean architecture.

## Features
- Product management
- Order processing
- User authentication
- File uploads
- Monitoring and logging
- Modular structure for easy maintenance

## Project Structure
- `main.go` — Entry point
- `modules/` — Business logic (products, orders, users, etc.)
- `pkg/` — Utilities, authentication, logging, database, etc.
- `config/` — Configuration files
- `assets/` — Static assets and logs
- `myTests/` — Test files

## Getting Started
1. **Clone the repository**
   ```powershell
   git clone https://github.com/tonrock01/another-world-shop.git
   cd another-world-shop
   ```
2. **Install Go (if not already installed)**
   - [Download Go](https://golang.org/dl/)
3. **Run the application**
   ```powershell
   go run main.go
   ```

## Configuration
- Edit settings in `config/config.go` as needed.

## Testing
- Run tests with:
   ```powershell
   go test ./myTests/...
   ```