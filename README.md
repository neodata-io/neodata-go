neodata-go/
├── config/                     # Configuration management
│   ├── config.go               # Wrapper around Viper for loading configurations
│   └── config.yaml             # Default config.yaml template
├── http/                       # HTTP and networking utilities
│   ├── client.go               # Pre-configured HTTP client with retries, timeouts
│   ├── middleware.go           # Common HTTP middleware (logging, rate limiting, etc.)
│   └── response.go             # Standardized JSON response formatting
├── logging/                    # Logging utilities
│   └── logger.go               # Centralized logging setup (Zap, Logrus, etc.)
├── database/                   # Database connection and helpers
│   ├── postgres.go             # Postgres connection pool setup
│   └── mongo.go                # MongoDB connection setup
├── caching/                    # Caching utilities
│   ├── redis.go                # Redis client setup and wrapper functions
│   └── memory.go               # In-memory caching utilities
├── messaging/                  # Messaging and event-driven utilities
│   ├── publisher.go            # Event publisher (Kafka, RabbitMQ, etc.)
│   └── subscriber.go           # Event subscriber
└── README.md                   # Documentation on using the library
