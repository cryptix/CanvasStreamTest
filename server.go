package main

import (
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/cryptix/canvas"
	"github.com/gorilla/mux"
)

func MakeGradientHandler(cv *canvas.Canvas, tick <-chan bool) websocket.Handler {

	return func(ws *websocket.Conn) {
		for {
			<-tick
			io.WriteString(ws, "NewImage")
		}
	}
}

func getImage(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "image/png")
	png.Encode(resp, cv)
}

var (
	world *World
	cv    *canvas.Canvas
)

func main() {
	tick := make(chan bool)

	width, height := 1024, 1024

	cv = canvas.NewCanvas(image.Rect(0, 0, width, height))

	n := 50

	go func() {
		for {
			tick <- true
			world = NewWorld(n, 6, cv)

			time.Sleep(time.Second * 5)
			n += 5
		}
	}()

	go func() {
		for {
			tick <- true
			world.SendMessages()

			time.Sleep(time.Second * 1)
		}
	}()

	router := mux.NewRouter()
	router.Handle("/ws", websocket.Handler(MakeGradientHandler(cv, tick)))

	router.HandleFunc("/getImage", getImage)

	staticHandler := http.FileServer(http.Dir("."))
	router.PathPrefix("/").Handler(staticHandler)

	log.Println("Listening on :3001")
	http.ListenAndServe("localhost:3001", router)
}
