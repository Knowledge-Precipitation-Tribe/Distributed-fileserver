package main

import (
	"Distributed-fileserver/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(":8080")
}
