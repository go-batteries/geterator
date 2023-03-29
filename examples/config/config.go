package config

//go:generate go run github.com/go-batteries/geterator -type=Config
type Config struct {
	Env   string
	Paths []string
	ID    int
}
