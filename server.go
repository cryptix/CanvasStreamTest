package main

import (
	"bytes"
	"encoding/base64"
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

func MakeGradientHandler(canvas *canvas.Canvas, tick <-chan bool) websocket.Handler {

	return func(ws *websocket.Conn) {
		for {

			<-tick

			imgBuf := new(bytes.Buffer)
			imgEncoder := base64.NewEncoder(base64.StdEncoding, imgBuf)

			// jpeg.Encode(imgEncoder, canvas, nil)
			png.Encode(imgEncoder, canvas)

			imgEncoder.Close()

			io.Copy(ws, imgBuf)

		}
	}
}

func main() {
	tick := make(chan bool)

	width, height := 640, 640

	canvas := canvas.NewCanvas(image.Rect(0, 0, width, height))

	n := 150
	world := NewWorld(n, 6, canvas)

	go func() {
		for {
			world = NewWorld(n, 6, canvas)
			tick <- true

			time.Sleep(time.Second * 5)

			n += 1
		}
	}()

	go func() {
		for {

			world.SendMessages(1)

			time.Sleep(250 * time.Millisecond)
			tick <- true

		}
	}()

	router := mux.NewRouter()
	router.Handle("/ws", websocket.Handler(MakeGradientHandler(canvas, tick)))

	staticHandler := http.FileServer(http.Dir("."))
	router.PathPrefix("/").Handler(staticHandler)

	log.Println("Listening on :3001")
	http.ListenAndServe("localhost:3001", router)
}
