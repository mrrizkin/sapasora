# Sapasora - GoFiber + Inertia.js + Vite Application Template

## Project Overview

Sapasora is a modern full-stack web application template that combines Go (with Fiber framework) for the backend and Vue.js with Inertia.js for the frontend. The project follows a modular architecture using Uber's FX dependency injection framework and includes modern tooling with Vite as the build system.

### Key Technologies & Features
- **Backend**: Go with Fiber framework (a fast HTTP framework)
- **Frontend**: Vue.js 3 with TypeScript, Inertia.js for server-side rendering
- **Build System**: Vite with hot module replacement
- **Styling**: Tailwind CSS with @nuxt/ui components
- **Database**: GORM with support for MySQL, PostgreSQL, SQLite, and other drivers
- **Dependency Injection**: Uber's FX framework for modular architecture
- **Templating**: `templ` package for type-safe Go templates
- **Code Quality**: ESLint, Prettier, TypeScript, Oxlint
- **Testing**: Vitest for unit tests, Playwright for E2E tests, Go testing for backend
- **Documentation**: Swagger API documentation
- **Session Management**: Multiple storage backends (memory, Redis, etc.)

### Architecture Components
- **GoFiber**: High-performance web framework written in Go
- **Inertia.js**: Client-side router that brings modern SPA responsiveness to server-side applications
- **Vite**: Next-generation frontend build tool with instant server start and lightning-fast hot updates
- **Templ**: Type-safe HTML templating for Go
- **FX**: Uber's dependency injection framework for modular Go applications
- **GORM**: Powerful ORM library for Go
- **Tailwind CSS**: Utility-first CSS framework

## Building and Running

### Prerequisites
- Go 1.25.1 or higher
- Node.js 20.0.0 or higher
- pnpm package manager

### Development Setup

```bash
# Clone the repository
git clone <repository-url>
cd sapasora

# Install dependencies
pnpm install

# Start development server with hot-reload
pnpm dev
```

The development server will start at `http://localhost:3000`, serving both the Go backend and Vite frontend with hot reloading.

### Alternative Development Commands

```bash
# Build the application (frontend and backend separately)
pnpm build

# Build for SSR (Server-Side Rendering)
pnpm build:ssr

# Start the Go application only
pnpm dev:app

# Start the Vite development server only
pnpm dev:assets

# Run SSR development server
pnpm dev:ssr
```

### Testing

```bash
# Run all tests (unit, E2E, and Go tests)
pnpm test

# Run unit tests only
pnpm test:unit

# Run E2E tests with Playwright
pnpm test:e2e

# Run Go tests only
pnpm test:go

# Run unit tests in watch mode
pnpm test:watch
```

### Code Quality

```bash
# Format code with Prettier
pnpm format

# Check formatting without applying changes
pnpm format:check

# Lint JavaScript/TypeScript code
pnpm lint

# Type check TypeScript
pnpm type-check

### Toolbox Commands

The application includes a comprehensive toolbox CLI for development tasks:

```bash
# Run the toolbox directly
go run ./cmd/toolbox

# Or use the pnpm alias
pnpm toolbox

# Generate autowired dependencies
go run ./cmd/toolbox autowired

# Generate scaffolding (models, contexts, JSON APIs, etc.)

# Scaffold a new model with GORM struct and basic methods
go run ./cmd/toolbox scaffold model module entity table_name [field_name:field_type...]
# Example: go run ./cmd/toolbox scaffold model user users name:string email:string

# Scaffold a context (repository pattern implementation)
go run ./cmd/toolbox scaffold context module entity table_name [field_name:field_type...]
# Example: go run ./cmd/toolbox scaffold context user users name:string email:string

# Scaffold a JSON API with full CRUD endpoints (includes model, repository interface/implementation, service, controller with list/get/store/update/delete methods, data transfer objects, and Swagger docs)
go run ./cmd/toolbox scaffold json module entity table_name [field_name:field_type...]
# Example: go run ./cmd/toolbox scaffold json user users name:string email:string

# Scaffold a policy for authorization (access control rules)
go run ./cmd/toolbox scaffold policy StructName policy_name
# Example: go run ./cmd/toolbox scaffold policy UserPolicy user
# Creates a policy file in internal/app/policies/ that implements authorization logic

# Scaffold authentication (currently incomplete - creates migration only)
go run ./cmd/toolbox scaffold auth module entity table_name [field_name:field_type...]
# Example: go run ./cmd/toolbox scaffold auth user users name:string email:string password:string
# Note: This command currently only creates a migration file and does not generate authentication implementation

# Database migrations
go run ./cmd/toolbox migrate run
go run ./cmd/toolbox migrate rollback
go run ./cmd/toolbox migrate reset
go run ./cmd/toolbox migrate fresh
go run ./cmd/toolbox migrate status
go run ./cmd/toolbox migrate make create_users_table

# Display routes
go run ./cmd/toolbox routes

# Generate Swagger documentation
go run ./cmd/toolbox swagger

# Generate TypeScript routes
go run ./cmd/toolbox wayfinder
```

## Project Structure

```
sapasora/
├── cmd/
│   ├── main/                 # Main application entry points
│   │   ├── main.go           # Main function entry point
│   │   └── serve.go          # Server startup and configuration
│   └── toolbox/              # CLI tools for scaffolding and management
│       ├── autowired.go      # Auto-generate dependency injection code
│       ├── kuproy.go         # Scaffolding commands (models, contexts, etc.)
│       ├── migrate.go        # Database migration management
│       ├── routes.go         # Display application routes
│       ├── swagger.go        # Generate OpenAPI documentation
│       ├── toolbox.go        # Main toolbox CLI entry point
│       └── wayfinder.go      # Generate TypeScript routes from Fiber routes
├── internal/
│   └── app/                  # Application bootstrap and core module
│       └── module.go         # Main application module using FX
├── platform/                 # Platform-specific modules and services
│   ├── cache/                # Cache implementations
│   ├── config/               # Configuration management
│   ├── database/             # Database configuration and connection
│   ├── logger/               # Logging utilities
│   ├── mail/                 # Email services
│   ├── oauth/                # OAuth integrations
│   ├── pubsub/               # Publish-subscribe messaging
│   ├── satpam/               # Authentication and authorization
│   ├── scheduler/            # Cron jobs and scheduled tasks
│   ├── server/               # HTTP server configuration
│   ├── session/              # Session management
│   ├── storage/              # File storage implementations
│   ├── support/              # Utility functions and helpers
│   ├── swagger/              # API documentation
│   ├── ui/                   # Frontend integration components
│   │   ├── bundler/          # Asset bundling
│   │   ├── inertia/          # Inertia.js integration
│   │   └── view/             # View rendering
│   └── validator/            # Input validation
├── resources/                # Frontend assets and templates
│   ├── assets/               # Static assets
│   ├── css/                  # CSS files (Tailwind)
│   ├── js/                   # JavaScript/TypeScript files
│   │   ├── app.ts            # Main Vue application entry
│   │   ├── ssr.ts            # Server-side rendering entry
│   │   └── colormode.ts      # Color mode management
│   └── views/                # Go template files (.templ)
├── public/                   # Public static files
├── storage/                  # Storage directory (logs, cache, etc.)
├── build/                    # Built application artifacts
├── .air.toml                 # Air configuration for Go hot-reload
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── package.json              # Node.js package manifest
├── pnpm-lock.yaml            # Dependency lock file
├── vite.config.ts            # Vite build configuration
├── tsconfig.json             # TypeScript configuration
└── README.md                 # Project documentation
```

## Development Conventions

### Backend (Go)
- Uses Uber's FX for dependency injection
- Modular architecture with separate platform packages
- Templatized Go code using `templ` package
- Structured logging with zerolog
- GORM for database operations
- Fiber middleware for common functionality

### Frontend (Vue.js + Inertia.js)
- TypeScript for type safety
- Composition API with `<script setup>` syntax
- Inertia.js for server-side rendering
- Tailwind CSS utility classes
- Nuxt UI components for consistent UI elements
- Vite for fast builds and hot module replacement

### Code Quality
- ESLint with Oxlint for JavaScript/TypeScript linting
- Prettier for code formatting
- TypeScript strict mode
- Go code formatted with gofmt
- Comprehensive test coverage encouraged

### Environment Variables
- Uses `.env` files with `godotenv` package
- Example `.env` values provided in `.env.example`

## Key Components

### Dependency Injection (FX Framework)
The application uses Uber's FX framework for dependency injection, allowing for modular and testable code. The main application module is bootstrapped in `internal/app/module.go` and orchestrates all platform modules.

### Inertia Integration
Inertia.js connects the Go backend with Vue.js frontend, providing a single-page application experience without building an API. The integration is handled in `platform/ui/inertia/`.

### Asset Bundling
The Vite build system handles asset bundling and development server. Laravel Vite Plugin is used to integrate with the Inertia workflow.

### Configuration
Configuration is managed through the `platform/config/` package, supporting environment variables and configuration files.

### Database Migrations
Database schema management is handled through GORM migrations in `platform/database/migrations/`.

### Toolbox CLI
The application includes a comprehensive toolbox CLI (`cmd/toolbox`) for development tasks:
- **Autowired**: Auto-generates dependency injection code for FX framework
- **Scaffolding**: Creates models, contexts, JSON APIs, policies, and authentication modules
  - **model**: Generates GORM model structs with fields and basic database operations, along with database migration files
  - **context**: Generates repository pattern implementation (interface and implementation) for database operations
  - **json**: Generates complete JSON API including models, repository, service, controllers with CRUD endpoints (list, get, store, update, delete), data transfer objects, and OpenAPI/Swagger documentation
  - **policy**: Generates authorization policies for access control rules
  - **auth**: (Currently incomplete) Creates a migration file for the authentication table but does not generate authentication implementation files
- **Migrations**: Manages database migrations (run, rollback, reset, fresh, make)
- **Routes**: Displays all registered application routes
- **Swagger**: Generates OpenAPI 3.0 documentation from code
- **Wayfinder**: Generates TypeScript routes from Fiber routes for frontend integration

## Deployment

### Production Build
```bash
# Create production build
pnpm build

# The build creates:
# - Compiled Go binary in ./build/app
# - Bundled frontend assets in ./public/
```

### Running in Production
After building, run the compiled application:
```bash
./build/app
```

The server will start on the configured port (defaults to 3000).
