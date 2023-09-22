package main

import "used-car-deal-gobackend/server"

func main() {

	server, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	server.Start()

}
