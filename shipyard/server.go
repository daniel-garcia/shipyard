package main

import "github.com/daniel-garcia/shipyard"


func main() {
	controller, _ := shipyard.NewServiceController("/tmp/test.json")
	controller.Save()
}

