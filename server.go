package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

type Message struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Heavily based on Kubernetes' (https://github.com/GoogleCloudPlatform/kubernetes) detection code.
var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Handle websockets if specified.
	if isWebsocketRequest(r) {
		websocket.Handler(HandleWebSockets).ServeHTTP(w, r)
	} else {
		HandleHttp(w, r)
	}
	glog.Info("Finished sending response...")
}

func HandleWebSockets(ws *websocket.Conn) {

	outputChan := make(chan string)
	go runBin(outputChan)

	for {
		select {
		case output, ok := <-outputChan:
			if !ok {
				// Channel closed, the runBin function has finished executing.
				glog.Info("runBin has finished executing.")
				return
			}
			// Print the output received from the channel.
			fmt.Print(output)
			m := Message{
				Id:      1,
				Message: output,
			}
			err := websocket.JSON.Send(ws, &m)
			if err != nil {
				glog.Infof("Client stopped listening...")
				return
			}
		default:
			// Add any other main loop functionality here, if needed.
			// For example, you can add a sleep to avoid busy-waiting.
			// time.Sleep(100 * time.Millisecond)
		}
	}


	// for i := 0; i < 5; i++ {
	// 	glog.Infof("Sending some data: %d", i)
	// 	m := Message{
	// 		Id:      i,
	// 		Message: fmt.Sprintf("Sending you \"%d\"", i),
	// 	}
	// 	err := websocket.JSON.Send(ws, &m)
	// 	if err != nil {
	// 		glog.Infof("Client stopped listening...")
	// 		return
	// 	}

	// 	// Artificially induce a 1s pause
	// 	time.Sleep(time.Second)
	// }
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.NotFound(w, r)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Send the initial headers saying we're gonna stream the response.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	enc := json.NewEncoder(w)

	for i := 0; i < 5; i++ {
		select {
		case <-cn.CloseNotify():
			glog.Infof("Client stopped listening")
			return
		default:
			// Artificially wait a second between reponses.
			time.Sleep(time.Second)

			glog.Infof("Sending some data: %d", i)
			m := Message{
				Id:      i,
				Message: fmt.Sprintf("Sending you \"%d\"", i),
			}

			// Send some data.
			err := enc.Encode(m)
			if err != nil {
				glog.Fatal(err)
			}
			flusher.Flush()
		}
	}
}

// runBin will send the output to the provided channel.
func runBin(outputChan chan<- string) {
	cmd := exec.Command("bash", "-c", "/workspace/main -m /workspace/llama-2-7b-chat.ggmlv3.q8_0.bin -n 1024")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1) // Read output byte by byte

	for {
		n, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // Reached end of output
			}
			log.Fatal(err)
		}
		outputChan <- string(buf[:n]) // Send the character read through the channel
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	close(outputChan) // Close the channel after the command finishes.
}

// Server.
func main() {
	flag.Parse()
	// outputChan := make(chan string)
	// go runBin(outputChan)

	// http.HandleFunc("/", Handle)

	// go func() {
	// 	for {
	// 		select {
	// 		case output, ok := <-outputChan:
	// 			if !ok {
	// 				// Channel closed, the runBin function has finished executing.
	// 				glog.Info("runBin has finished executing.")
	// 				return
	// 			}
	// 			// Print the output received from the channel.
	// 			fmt.Print(output)
	// 		default:
	// 			// Add any other main loop functionality here, if needed.
	// 			// For example, you can add a sleep to avoid busy-waiting.
	// 			// time.Sleep(100 * time.Millisecond)
	// 		}
	// 	}
	// }()

	glog.Infof("Serving...")
	http.HandleFunc("/", Handle)

	// glog.Infof("Serving...")
	glog.Fatal(http.ListenAndServe(":8080", nil))

	// for {
	// 	select {
	// 	case output, ok := <-outputChan:
	// 		if !ok {
	// 			// Channel closed, the runBin function has finished executing.
	// 			glog.Info("runBin has finished executing.")
	// 			return
	// 		}
	// 		// Print the output received from the channel.
	// 		fmt.Print(output)
	// 	default:
	// 		// Add any other main loop functionality here, if needed.
	// 		// For example, you can add a sleep to avoid busy-waiting.
	// 		// time.Sleep(100 * time.Millisecond)
	// 	}
	// }

}
