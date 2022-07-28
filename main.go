package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
	"wxcloudrun-golang/db"
)

func main() {
	var try = 3
	var start = func() error {
		if err := db.Init(); err != nil {
			log.Println(fmt.Sprintf("mysql init failed with %+v", err))
			try--
			return err
		}
		return nil
	}
	for try > 0 {
		err := start()
		if err != nil {
			time.Sleep(time.Second)
		}
		break
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
