package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"wxcloudrun-golang/db"
)

func main() {
	if err := db.Init(); err != nil {
		log.Println(fmt.Sprintf("mysql init failed with %+v", err))
	}
	router := gin.Default()
	myRouter(router)
	port := os.Getenv("MY_PORT")
	if port == "" {
		port = ":8889"
	} else {
		port = ":" + port
	}
	log.Fatalln(router.Run(port))
}
