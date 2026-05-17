package main

import (
	"configurator/cmd/config"
	"context"
)

func main() {
	config.Execute(context.Background())
}
