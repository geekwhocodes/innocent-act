package main

// Config /
type Config struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
	Db   struct {
		Provider string `koanf:"provider"`
		Sqlite   struct {
			Name string `koanf:"name"`
			Path string `koanf:"path"`
		} `koanf:"sqlite"`
	} `koanf:"db"`
}
