package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	r := gin.Default()

	r.GET("/ws", createWS)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// var homeTemplate = template.Must(template.New("").Parse(`
// <!DOCTYPE html>
// <html>
// <head>
// <meta charset="utf-8">
// <script>
// window.addEventListener("load", function(evt) {
//     var output = document.getElementById("output");
//     var input = document.getElementById("input");
//     var ws;
//     var print = function(message) {
//         var d = document.createElement("div");
//         d.textContent = message;
//         output.appendChild(d);
//         output.scroll(0, output.scrollHeight);
//     };
//     document.getElementById("open").onclick = function(evt) {
//         if (ws) {
//             return false;
//         }
//         ws = new WebSocket("{{.}}");
//         ws.onopen = function(evt) {
//             print("OPEN");
//         }
//         ws.onclose = function(evt) {
//             print("CLOSE");
//             ws = null;
//         }
//         ws.onmessage = function(evt) {
//             print("RESPONSE: " + evt.data);
//         }
//         ws.onerror = function(evt) {
//             print("ERROR: " + evt.data);
//         }
//         return false;
//     };
//     document.getElementById("send").onclick = function(evt) {
//         if (!ws) {
//             return false;
//         }
//         print("SEND: " + input.value);
//         ws.send(input.value);
//         return false;
//     };
//     document.getElementById("close").onclick = function(evt) {
//         if (!ws) {
//             return false;
//         }
//         ws.close();
//         return false;
//     };
// });
// </script>
// </head>
// <body>
// <table>
// <tr><td valign="top" width="50%">
// <p>Click "Open" to create a connection to the server,
// "Send" to send a message to the server and "Close" to close the connection.
// You can change the message and send multiple times.
// <p>
// <form>
// <button id="open">Open</button>
// <button id="close">Close</button>
// <p><input id="input" type="text" value="Hello world!">
// <button id="send">Send</button>
// </form>
// </td><td valign="top" width="50%">
// <div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
// </td></tr></table>
// </body>
// </html>
// `))
