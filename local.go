package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/alexgiesting/gillings-search/paths"
	"github.com/alexgiesting/gillings-search/query"
	"github.com/alexgiesting/gillings-search/update"
)

func runMongod() *os.Process {
	os.Setenv(paths.ENV_MONGODB_HOST, "127.0.0.1")
	os.Setenv(paths.ENV_MONGODB_PORT, "27017")
	os.Setenv(paths.ENV_MONGODB_NAME, "test")
	os.Unsetenv(paths.ENV_MONGODB_ADMIN_PASSWORD)

	MONGOD_PATH, err := exec.LookPath("mongod")
	if err != nil {
		log.Fatal(err)
	}
	MONGOD_ARGV := strings.Split("mongod --config ./mongodb/mongod.cfg", " ")
	procAttr := os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	mongod, err := os.StartProcess(MONGOD_PATH, MONGOD_ARGV, &procAttr)
	if err != nil {
		log.Fatal(err)
	}

	return mongod
}

func loadKey(envVariable string, keyFilename string) {
	keyFile, err := os.Open(keyFilename)
	if err != nil {
		log.Fatal(err)
	}
	apiKey, err := io.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv(envVariable, string(apiKey))
	keyFile.Close()
}

const (
	SERVER_PORT = ":3000"
	QUERY_PORT  = ":3001"
	UPDATE_PORT = ":3002"
)

func proxyAtPort(serveMux *http.ServeMux, path string, port string) {
	url, err := url.Parse(fmt.Sprintf("http://localhost%s/", port))
	if err != nil {
		log.Fatal(err)
	}
	serveMux.Handle(path, httputil.NewSingleHostReverseProxy(url))
}

func main() {
	mongod := runMongod()
	defer mongod.Kill()

	os.Setenv(paths.ENV_QUERY_PORT, QUERY_PORT)
	os.Setenv(paths.ENV_UPDATE_PORT, UPDATE_PORT)
	loadKey(paths.ENV_SCOPUS_API_KEY, "scopus.key")
	loadKey(paths.ENV_SCOPUS_CLIENT_ADDRESS, "subscriber.key")
	loadKey(paths.ENV_UPDATE_KEY, "update.key")

	// TODO prefix logs with process name, so we can tell them apart
	//      maybe use contexts?
	go query.Main()
	go update.Main()

	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir("./static")))
	proxyAtPort(serveMux, paths.PATH_QUERY, os.Getenv(paths.ENV_QUERY_PORT))
	proxyAtPort(serveMux, paths.PATH_UPDATE, os.Getenv(paths.ENV_UPDATE_PORT))
	log.Fatal(http.ListenAndServe(SERVER_PORT, serveMux))
}
