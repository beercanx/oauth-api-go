package server

type Config struct {
	HttpHost string `env:"HTTP_HOST" default:"0.0.0.0"`
	HttpPort string `env:"HTTP_PORT" default:"8080"`
}
