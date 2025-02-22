package server

type Config struct {
	HttpHost string `env:"HTTP_HOST" default:"127.0.0.1"`
	HttpPort string `env:"HTTP_PORT" default:"8080"`
}
