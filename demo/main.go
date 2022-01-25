package main

import (
	"fmt"
	"net/http"

	"github.com/shara/helix/config"
	"github.com/shara/helix/services"
)

func helloworldHandler(rw http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(rw, "Hello World")
}

func main() {
	app := services.NewApplication("demo", 8888, "v1", config.DATABASE_DNS)
	//Hello World

	app.HandleFunc("/v1/helloworld", helloworldHandler)

	app.RunService()

}
