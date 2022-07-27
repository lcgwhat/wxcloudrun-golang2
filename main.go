package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"wxcloudrun-golang/db"
)

func main() {
	if err := db.Init(); err != nil {
		log.Println(fmt.Sprintf("mysql init failed with %+v", err))
	}
	router := gin.Default()
	myRouter(router)

	log.Fatalln(router.Run(":80"))
}
