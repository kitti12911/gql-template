# gql-template

GraphQL gateway using [gqlgen](https://gqlgen.com). sits in front of gRPC services and exposes a unified GraphQL API. currently connects to [grpc-template](https://github.com/kitti12911/grpc-template)'s example service.

## features

- GraphQL gateway with schema-per-service layout
- connects to gRPC backend services (example service)
- opentelemetry tracing for operations, responses, and field resolvers (exports to alloy/tempo via OTLP)
- `@auth` directive (print-only, keycloak integration planned)
- structured logging with slog
- config from yaml or environment variables
- GraphQL playground at `/`
- graceful shutdown
- hot reload with air

## requirements

- go 1.26.0 or higher
- [buf](https://buf.build/) installed
- protoc-gen-go installed
- protoc-gen-go-grpc installed
- a running gRPC backend (e.g. [grpc-template](https://github.com/kitti12911/grpc-template))

## optional

- [air](https://github.com/air-verse/air) for hot reload

## setup

### install go

- macos:

    ```bash
    brew install go
    ```

- linux (apt):

    ```bash
    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt update
    sudo apt install golang-go
    ```

- linux (snap):

    ```bash
    sudo snap install --classic go
    ```

then add go bin to your PATH:

- macos (zsh):

    ```bash
    echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
    source ~/.zshrc
    ```

- linux (bash):

    ```bash
    echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
    source ~/.bashrc
    ```

## install buf

- macos:

    ```bash
    brew install bufbuild/buf/buf
    ```

- linux:

    ```bash
    # see https://buf.build/docs/cli/installation for other methods
    curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m) -o /usr/local/bin/buf
    chmod +x /usr/local/bin/buf
    ```

### install protoc plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1
```

### install air (optional)

```bash
go install github.com/air-verse/air@v1.64.5
```

## project structure

```bash
gql-template/
├── cmd/
│   └── server/
│       └── main.go                    # entrypoint
├── internal/
│   ├── config/
│   │   └── config.go                  # config struct
│   ├── directive/
│   │   └── auth.go                    # @auth directive (print only)
│   └── server/
│       ├── http.go                    # HTTP server + gqlgen handler
│       └── otel.go                    # opentelemetry gqlgen extension
├── graph/
│   ├── generated.go                   # gqlgen generated (do not edit)
│   ├── model/
│   │   └── models_gen.go              # gqlgen generated models (do not edit)
│   ├── resolver.go                    # root resolver with gRPC clients
│   ├── base.resolvers.go             # gqlgen generated (do not edit)
│   ├── example.resolvers.go          # example service resolvers
│   └── schema/
│       ├── base.graphqls              # shared directives, scalars, root types
│       └── example.graphqls           # example service schema
├── gen/                               # generated protobuf code (do not edit)
├── .air.toml                          # air config
├── buf.gen.yaml                       # buf code generation config
├── config.yml                         # app config (gitignored)
├── Dockerfile
├── gqlgen.yml                         # gqlgen config
├── Makefile
└── go.mod
```

schemas and resolvers are split per service. each service gets its own `.graphqls` schema file in `graph/schema/` and gqlgen generates a matching resolver file in `graph/`.

## config

create a `config.yml` in the project root:

```yaml
service_name: gql-template
port: 8080
collector_endpoint: localhost
collector_port: 4317
shutdown_timeout: 10s

logging:
  level: info
  add_source: false
  service_name: gql-template
  enable_trace: true

example_service:
  host: localhost
  port: 50051
```

config values can be set via environment variables:

| yaml key               | env variable         | description                   |
|------------------------|----------------------|-------------------------------|
| `service_name`         | `SERVICE_NAME`       | service name for tracing/logs |
| `port`                 | `PORT`               | HTTP server port              |
| `collector_endpoint`   | `COLLECTOR_ENDPOINT` | OTLP collector host           |
| `collector_port`       | `COLLECTOR_PORT`     | OTLP collector port           |
| `shutdown_timeout`     | `SHUTDOWN_TIMEOUT`   | graceful shutdown timeout     |
| `logging.level`        | `LOG_LEVEL`          | debug, info, warn, error      |
| `logging.add_source`   | `LOG_ADD_SOURCE`     | include source file in logs   |
| `logging.service_name` | `LOG_SERVICE_NAME`   | service name in log output    |
| `logging.enable_trace` | `LOG_ENABLE_TRACE`   | add trace/span id to logs     |

## how to run

with air (hot reload):

```bash
make air
```

without air:

```bash
make run
```

then open <http://localhost:8080> for the GraphQL playground.

## available commands

| command            | description                              |
|--------------------|------------------------------------------|
| `make air`         | run with hot reload                      |
| `make run`         | run without hot reload                   |
| `make tidy`        | go mod tidy                              |
| `make fmt`         | format code                              |
| `make test`        | run tests with race detector             |
| `make cov`         | run tests with coverage report           |
| `make gen`         | regenerate gqlgen + protobuf code        |
| `make gen-gql`     | regenerate gqlgen code only              |
| `make gen-proto`   | regenerate protobuf code only            |

## generate code

gqlgen code (models, resolvers):

```bash
make gen-gql
```

protobuf client stubs from [proto-template](https://github.com/kitti12911/proto-template):

```bash
make gen-proto
```

both:

```bash
make gen
```

## adding a new service

1. add proto definitions to the [proto-template](https://github.com/kitti12911/proto-template) repo
2. add the path to `make gen-proto` in `Makefile` and `buf.gen.yaml` M-flag mappings
3. run `make gen-proto`
4. create a new schema file `graph/schema/your-service.graphqls`:

    ```graphql
    extend type Query {
      yourThing(id: ID!): YourThing!
    }

    type YourThing {
      id: ID!
      name: String!
    }
    ```

5. run `make gen-gql` — gqlgen creates `graph/your-service.resolvers.go`
6. implement the resolver methods in the generated file
7. add the gRPC client to `graph/resolver.go` and wire it up in `cmd/server/main.go`
8. add the service address to `config.yml`
