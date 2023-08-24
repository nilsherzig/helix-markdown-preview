package main

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"embed"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//go:embed static
var embStaticFS embed.FS

//go:embed templates
var emdTemplateFS embed.FS

func main() {

	// golang error handling meme
	var err error

	// choose root folder to monitor - use current folder if no args are provided
	var folderPath string

	if len(os.Args) != 2 {
		folderPath, err = os.Getwd()
		handleErr(err)
		log.Printf("No path specified, using %s\n", folderPath)
	} else {
		folderPath = os.Args[1]
		folderPath = filepath.Clean(folderPath)
		log.Println("working in:", folderPath)
	}

	// channel logic, used by filewatcher goroutine to exchange updates with webserver
	channel := make(chan string)

	message := "noch keine änderungen"
	go func(channel chan string) {
		for {
			message = <-channel
		}
	}(channel)

	// webserver setup
	// debug build reads files on every site refresh from filesystem
	// prod build embeds files into binary
	// different logging
	// compile using `DEBUG=true`
	var router *gin.Engine

	debug := os.Getenv("DEBUG")
	if len(debug) > 0 {
		router = gin.Default()
		router.LoadHTMLGlob("templates/*")
		router.Static("/static", "./static")
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		templ := template.Must(template.New("").ParseFS(emdTemplateFS, "templates/*.html"))
		router.SetHTMLTemplate(templ)

		sub, err := fs.Sub(embStaticFS, "static")
		handleErr(err)
		router.StaticFS("/static", http.FS(sub)) // embed
	}

	// index endpoint
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"folder": folderPath,
		})
	})

	// websocket endpoint
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
				err := conn.WriteMessage(websocket.TextMessage, []byte(currentMessage))
				handleErr(err)
			}
			oldMessage = currentMessage
			time.Sleep(time.Millisecond * 300) // there has to be another way TODO
		}
	})

	// start file watcher in goroutine
	go watchFolder(folderPath, channel)

	// find open port, starting from 8080
	webserverPort := 8080
	for {
		_, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", webserverPort), time.Second)
		if err != nil {
			break
		}
		webserverPort += 1
	}

	// start gin webserver
	log.Printf("website running at http://127.0.0.1:%d\n", webserverPort)
	openbrowser(fmt.Sprintf("http://127.0.0.1:%d", webserverPort))
	err = router.Run(fmt.Sprintf(":%d", webserverPort))
	handleErr(err)
}
