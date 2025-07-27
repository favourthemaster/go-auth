# Authentication Service

This project is a Go-based authentication service using Fiber, designed for secure user management and session handling. It demonstrates best practices in structuring Go applications, handling authentication, and integrating with databases.

## Features
- User registration and login
- JWT and session-based authentication
- Password hashing and validation
- Email notifications
- Redis integration for session management
- Modular code structure

## Setup Instructions
1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd authentication
   ```
2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **Configure environment:**
   - Set up your database and Redis connection in the config files.
   - Add any required environment variables.
4. **Run the application:**
   ```sh
   go run src/cmd/main.go
   ```

## Testing
Run unit tests with:
```sh
go test ./src/...
```

## Project Structure
- `src/cmd/` - Entry point
- `src/internal/` - Core logic (auth, db, user)
- `src/dto/` - Data transfer objects
- `src/errs/` - Error handling
- `src/models/` - Data models
- `src/utils/` - Utility functions

## License
See [LICENSE](LICENSE).

