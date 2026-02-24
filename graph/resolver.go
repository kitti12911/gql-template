package graph

import examplev1 "gql-template/gen/grpc/example/v1"

type Resolver struct {
	ExampleClient examplev1.ExampleServiceClient
}
