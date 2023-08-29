package main

import "github.com/mikalai2006/geoinfo-api/internal/app"

func main() {
	// base path for config: default = ./ (for test ../)
	const configPath = "./"

	app.Run(configPath)
}
