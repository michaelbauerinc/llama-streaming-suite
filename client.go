package main

import (
        "encoding/json"
        "flag"
        "io"
        "net/http"
        "fmt"
        "github.com/golang/glog"
        "golang.org/x/net/websocket"
)

var useWebsockets = flag.Bool("websockets", true, "Whether to use websockets")

type Message struct {
        Id      int    `json:"id,omitempty"`
        Message string `json:"message,omitempty"`
}

// Client.
func main() {
        flag.Parse()

        if *useWebsockets {
                ws, err := websocket.Dial("ws://0.0.0.0:8080/", "", "http://0.0.0.0:8080")
                for {
                        var m Message
                        err = websocket.JSON.Receive(ws, &m)
                        if err != nil {
                                if err == io.EOF {
                                        break
                                }
                                glog.Fatal(err)
                        }
                        fmt.Println("Received: %+v", m.Message)
                }
        } else {
                glog.Info("Sending request...")
                req, err := http.NewRequest("GET", "http://0.0.0.0:8080", nil)
                if err != nil {
                        glog.Fatal(err)
                }
                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                        glog.Fatal(err)
                }
                if resp.StatusCode != http.StatusOK {
                        glog.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
                }

                dec := json.NewDecoder(resp.Body)
                for {
                        var m Message
                        err := dec.Decode(&m)
                        if err != nil {
                                if err == io.EOF {
                                        break
                                }
                                glog.Fatal(err)
                        }
                        fmt.Println("Got response: %+v", m)
                }
        }

        fmt.Println("Server finished request...")
}
