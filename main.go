package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//if err := db.Init(); err != nil {
	//	panic(fmt.Sprintf("mysql init failed with %+v", err))
	//}
	router := gin.Default()
	myRouter(router)

	log.Fatalln(router.Run(":80"))
}
