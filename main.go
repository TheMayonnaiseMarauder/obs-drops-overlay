package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	//go:embed control.html
	controlHTML string
	//go:embed overlay.html
	overlayHTML       string
	websocketUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	controlChan       chan []byte
	statusChan        chan []byte
	localIP           string
	port              string
	mediaExt          = []string{".mp4", ".mkv", ".mp3"}
)

type FileListResponse struct {
	Files []string `json:"files"`
}

type Template struct {
	IP   string
	Port string
}

func listDir(path string) (files []string, err error) {
	var pathContents []os.DirEntry
	if pathContents, err = os.ReadDir(path); err != nil {
		return
	} else {
		for _, entry := range pathContents {
			for _, ext := range mediaExt {
				if filepath.Ext(entry.Name()) == ext {
					files = append(files, entry.Name())
					break
				}
			}
		}
	}
	return
}

func applyTemplate(htmlString string, w http.ResponseWriter) error {
	t := template.New("t")
	if _, err := t.Parse(htmlString); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("ERROR::t.Parse(rootHtml)::%v\n", err)
	}
	if err := t.Execute(w, Template{IP: localIP, Port: port}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("ERROR::t.Execute(w, &t)::%v\n", err)
	}
	return nil
}

func start() {
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		rType := r.FormValue("type")
		music, err := listDir(fmt.Sprintf("assets/%s/", rType))
		if err != nil {
			log.Panicf("listDir(\"assets/%s/\"):%v\n", rType, err)
		}
		bytes, err := json.Marshal(FileListResponse{Files: music})
		if err != nil {
			log.Panicf("json.Marshal(FileListResponse{Files: %s}):%v\n", rType, err)
		}
		if _, err := w.Write(bytes); err != nil {
			log.Panicf("w.Write(bytes):%v\n", err)
		}
	})
	http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		http.ServeFile(w, r, fmt.Sprintf("assets/%s", r.FormValue("file")))
	})
	http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
		if err := applyTemplate(controlHTML, w); err != nil {
			log.Fatalf("%v", err)
		}
	})
	http.HandleFunc("/overlay", func(w http.ResponseWriter, r *http.Request) {
		if err := applyTemplate(overlayHTML, w); err != nil {
			log.Fatalf("%v", err)
		}
	})
	http.HandleFunc("/overlayWS", func(w http.ResponseWriter, r *http.Request) {
		websocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("%v\n", err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing /overlayWS connection: %v\n", err)
			}
		}()
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				statusChan <- msg
			}
		}()
		for {
			if err = conn.WriteMessage(1, <-controlChan); err != nil {
				return
			}
		}
	})
	http.HandleFunc("/controlWS", func(w http.ResponseWriter, r *http.Request) {
		websocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("%v\n", err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing /controlWS connection: %v\n", err)
			}
		}()
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				controlChan <- msg
			}
		}()
		for {
			if err = conn.WriteMessage(1, <-statusChan); err != nil {
				return
			}
		}
	})
	fmt.Printf("starting server @ http://%s:%s/control http://%s:%s/overlay\n", localIP, port, localIP, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("%v\n", err)
	}
}

func init() {
	localIP = func() string {
		adders, err := net.InterfaceAddrs()
		if err != nil {
			return ""
		}
		for _, address := range adders {
			if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
				if inet.IP.To4() != nil {
					return inet.IP.String()
				}
			}
		}
		return ""
	}()
}

func main() {
	port = "8605"
	controlChan, statusChan = make(chan []byte), make(chan []byte)
	start()
}
