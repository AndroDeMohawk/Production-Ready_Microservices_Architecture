module github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment

go 1.26.0

require (
	github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared v0.0.0-00010101000000-000000000000
	github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.44.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared => ./../shared
replace github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform => ./../platform
