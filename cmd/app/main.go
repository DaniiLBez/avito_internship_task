package main

import "github.com/DaniiLBez/avito_internship_task/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
