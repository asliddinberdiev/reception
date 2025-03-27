package main

import "github.com/asliddinberdiev/reception/internal/app"

const confDir = "config"

func main() {
	app.Run(confDir)
}
