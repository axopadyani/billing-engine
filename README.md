# Billing Engine

## Overview

Billing Engine is a Go-based gRPC service that manages loan creation, retrieval, and payment processing. It provides a robust backend for handling financial transactions related to loans.

## Features

- Create new loans for users
- Retrieve current loan details
- Process loan payments
- gRPC API for efficient communication

## Prerequisites

- Docker
- Go 1.23.3

## Installation & Running

1. Clone the repository:

    ```shell
    git clone https://github.com/axopadyani/billing-engine.git
    cd billing-engine
    ```

2. Download Go module dependencies:

    ```shell
    go mod download
    ```

## Usage

1. Spin up needed dependencies:

    ```shell
    docker compose up
    ```

2. Set up the environment variables (update the `.env` file as needed):

    ```shell
    cp .env.sample .env
    ```

3. Run the server:

    ```shell
    go run cmd/server/main.go
    ```

## API

The service exposes the following gRPC methods:

- `CreateLoan`: Create a new loan for a user
- `GetCurrentLoan`: Retrieve the current loan details for a user
- `MakePayment`: Process a payment for a specific loan

For detailed API documentation, refer to the proto files in the `proto/v1` directory.

## Development

The tools needed for development:
- Protocol Buffers compiler (protoc)
- `mockgen` tool: https://github.com/uber-go/mock

To regenerate gRPC code when changes are made to the protobuf contract:

```shell
make gen-proto
```

To regenerate interface mocks when changes are made to the interfaces:
```shell
go generate ./...
```
