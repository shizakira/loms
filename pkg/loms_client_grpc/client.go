package loms_client_grpc

type Config struct {
	Host string `default:"localhost" envconfig:"GRPC_CLIENT_HOST"`
	Port string `default:"50051"     envconfig:"GRPC_CLIENT_PORT"`
}
