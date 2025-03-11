# Docker-Based Testing for Git Manager

This document explains how to run the Git Manager test suite in Docker containers, providing a consistent and isolated environment for testing.

## Prerequisites

- Docker installed and running
- Docker Compose installed (optional, for Docker Compose-based testing)
- Task installed (optional, for running Task-based commands)

## Running Tests in Docker

### Using Task (Recommended)

The simplest way to run tests in Docker is using Task:

```bash
# Run specific tests in Docker
task test:docker:specific -- ./internal/worktree

# Run tests with coverage in Docker
task test:docker:cover

# Run tests using Docker Compose
task test:docker:compose
```

### Using Scripts Directly

You can also run the scripts directly:

```bash
# Run all tests in Docker
./scripts/run-tests.sh

# Run specific tests in Docker
./scripts/run-tests.sh ./internal/worktree

# Run tests using Docker Compose
./scripts/run-tests-compose.sh

# Run specific tests using Docker Compose
./scripts/run-tests-compose.sh ./internal/worktree
```

### Using Docker Compose Directly

You can also use Docker Compose directly:

```bash
# Run all tests
docker-compose -f docker-compose.test.yml up --build test

# Run specific tests
docker-compose -f docker-compose.test.yml run --rm test-specific go test -v ./internal/worktree
```

## Test Results

Test results, including coverage reports, are stored in the `test-results` directory:

- `test-results/coverage.out`: Coverage report in Go's coverage format
- `test-results/coverage.html`: HTML coverage report (generated when running `task test:docker:cover`)

## Customizing the Docker Environment

### Dockerfile.test

The `Dockerfile.test` file defines the Docker image used for testing. You can modify this file to:

- Change the Go version
- Install additional dependencies
- Configure the git environment

### docker-compose.test.yml

The `docker-compose.test.yml` file defines the Docker Compose services for testing. You can modify this file to:

- Add additional services (e.g., databases)
- Configure environment variables
- Mount additional volumes

## Continuous Integration

This Docker-based testing setup is ideal for CI/CD pipelines. Here's an example GitHub Actions workflow:

```yaml
name: Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker
        uses: docker/setup-buildx-action@v2
      
      - name: Run tests in Docker
        run: ./scripts/run-tests.sh
      
      - name: Upload coverage report
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: test-results/
```

## Troubleshooting

### Tests Fail in Docker but Pass Locally

This could be due to:

1. Different Git version in Docker
2. Different filesystem behavior
3. Environment variables not set in Docker

Check the Docker logs for more information:

```bash
docker-compose -f docker-compose.test.yml logs test
```

### Docker Build Fails

If the Docker build fails, check:

1. Docker is running
2. You have sufficient permissions
3. The Go version in the Dockerfile is available

### Tests Hang or Take Too Long

Git operations can sometimes be slow in Docker. Try:

1. Increasing the timeout in your tests
2. Using a volume mount for the git repository to improve performance 
