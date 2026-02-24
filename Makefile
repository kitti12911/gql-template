# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server/main.go

fmt:
	go fmt ./...

test:
	env CGO_ENABLE=1 go test --race -v ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ____________________ Generate Command ____________________
gen: gen-nullable gen-gql gen-proto

gen-gql:
	go tool gqlgen generate

gen-nullable:
	go run ./cmd/gen-nullable

gen-proto:
	rm -rf gen/grpc
	buf generate https://github.com/kitti12911/proto-template.git --path common/v1 --path example/v1
