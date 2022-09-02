module gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs

go 1.19

replace gitlab.internal.cloud.payly.com.br/microservices/chassis/events => gitlab.internal.cloud.payly.com.br/microservices/chassis/events.git v0.0.1

replace gitlab.internal.cloud.payly.com.br/microservices/chassis/logging => gitlab.internal.cloud.payly.com.br/microservices/chassis/logging.git v0.0.3

require (
	gitlab.internal.cloud.payly.com.br/microservices/chassis/events v0.0.0-00010101000000-000000000000
	gitlab.internal.cloud.payly.com.br/microservices/chassis/logging v0.0.0-00010101000000-000000000000
	go.uber.org/dig v1.15.0
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)
