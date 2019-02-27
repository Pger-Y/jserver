package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"
)

func debug(data []byte, err error) {
	if err == nil {
		log.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

type Server struct {
	srv  *http.Server
	done chan struct{}
	data map[string]interface{}
}

type HT struct {
	data         map[string]interface{}
	root_content string
}

func (ht *HT) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	debug(httputil.DumpRequest(req, true))
	path := req.URL.Path
	if req.Method == "GET" {
		path = fmt.Sprintf("%s?%s", path, req.URL.RawQuery)
	}
	data, ok := ht.data[path]
	if ok {
		json.NewEncoder(w).Encode(data)
	} else if path == "/" {
		fmt.Fprintf(w, ht.root_content)
	} else {
		message := fmt.Sprintf("[%s] 404 Not Found", path)
		http.Error(w, message, http.StatusNotFound)
	}
}

func NewServer(addr string, data map[string]interface{}) *Server {
	paths := make([]string, 0, len(data))
	for k, _ := range data {
		paths = append(paths, k)
	}

	sort.Strings(paths)
	ps := strings.Join(paths, "\n")
	content := fmt.Sprintf("Supported path ==>\n%s\n", ps)

	ht := &HT{data: data, root_content: content}
	srv := &http.Server{Addr: addr, Handler: ht}
	s := &Server{srv: srv}
	s.done = make(chan struct{})

	return s
}

func (s *Server) Stop() {
	if s != nil {
		if err := s.srv.Shutdown(nil); err != nil {
			log.Println("Server Stop error,", err)
		}
		log.Println("Server stoped...")
		<-s.done
	}
}

func (s *Server) Start() {
	log.Println("Starting server...")
	go func() {
		s.srv.RegisterOnShutdown(func() {
			if s != nil {
				close(s.done)
			}
		})
		if err := s.srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver ListenAndServe() error:%s\n", err)
		}
	}()
}
