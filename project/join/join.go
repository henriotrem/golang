package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/negroni"
)

var (
	config   configuration
	upgrader = websocket.Upgrader{} // use default options
)

type configuration struct {
	Debug         bool   `default:"true"`
	Scheme        string `default:"http"`
	ListenAddress string `default:":8010"`
	PrivateKey    string `default:"ssl/server.key"`
	Certificate   string `default:"ssl/server.pem"`
}

type httpErr struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func handleErr(w http.ResponseWriter, err error, status int) {
	msg, err := json.Marshal(&httpErr{
		Msg:  err.Error(),
		Code: status,
	})
	if err != nil {
		msg = []byte(err.Error())
	}
	http.Error(w, string(msg), status)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
	defer c.Close()
	network := CreateNetwork(c)
	defer network.Quit()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			handleErr(w, err, http.StatusInternalServerError)
			break
		}
		if mt != websocket.TextMessage {
			handleErr(w, errors.New("only text message are supported"), http.StatusNotImplemented)
			break
		}
		handleNetworkRequest(c, network, string(msg))
	}
}

func handleNetworkRequest(c *websocket.Conn, network *Network, cmd string) {

	switch cmd {
	case "INIT":
		network.Init()
	case "NEW":
		network.New()
	case "PRINT":
		network.Print()
	case "QUIT":
		network.Quit()
		network = CreateNetwork(c)
	}
}

func main() {

	// Default values
	err := envconfig.Process("SOCKETCAM", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	if config.Debug {
		log.Printf("==> SCHEME: %v", config.Scheme)
		log.Printf("==> ADDRESS: %v", config.ListenAddress)
		log.Printf("==> PRIVATEKEY: %v", config.PrivateKey)
		log.Printf("==> CERTIFICATE: %v", config.Certificate)
	}

	router := newRouter()
	n := negroni.Classic()

	n.UseHandler(router)
	if config.Scheme == "https" {
		log.Fatal(http.ListenAndServeTLS(config.ListenAddress, config.Certificate, config.PrivateKey, n))

	} else {
		log.Fatal(http.ListenAndServe(config.ListenAddress, n))

	}
}

// NewRouter is the constructor for all my routes
func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/ws").
		Name("Communication Channel").
		HandlerFunc(serveWs)

	router.
		Methods("GET").
		PathPrefix("/").
		Name("Static").
		Handler(http.FileServer(http.Dir("./htdocs")))
	return router
}
