package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Please provide a path")
	}

	folderPath := os.Args[1]

	folderPath = filepath.Clean(folderPath)

	log.Println("working in:", folderPath)

	webserverPort := ":8080"

	channel := make(chan string)
	go watchFolder(folderPath, channel)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	message := "noch keine änderungen"

	go func(channel chan string) {
		for {
			message = <-channel
		}
	}(channel)

	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		oldMessage := message
		for {
			currentMessage := message
			if currentMessage != oldMessage {
				conn.WriteMessage(websocket.TextMessage, []byte(currentMessage))
			}
			oldMessage = currentMessage
			time.Sleep(time.Millisecond * 300)
		}
	})

	log.Printf("http://127.0.0.1%s\n", webserverPort)
	router.Run(webserverPort)

}

func exit(i int) {
	panic("unimplemented")
}
