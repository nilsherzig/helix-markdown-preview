package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//go:embed static
var embStaticFS embed.FS

//go:embed templates
var emdTemplateFS embed.FS

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
	// router.LoadHTMLGlob()
	// router.LoadHTMLFiles("./templates/index.html")

	templ := template.Must(template.New("").ParseFS(emdTemplateFS, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	// router.Static("/static", "./static")
	sub, _ := fs.Sub(embStaticFS, "static") // needs eror checks TODO

	router.StaticFS("/static", http.FS(sub)) // embed

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"folder": folderPath,
		})
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
			time.Sleep(time.Millisecond * 300) // there has to be another way TODO
		}
	})

	log.Printf("http://127.0.0.1%s\n", webserverPort)
	router.Run(webserverPort)

}

func exit(i int) {
	panic("unimplemented")
}
