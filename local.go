package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/alexgiesting/gillings-search/paths"
	"github.com/alexgiesting/gillings-search/poll"
	"github.com/alexgiesting/gillings-search/query"
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

func setupScopus() {
	keyFile, err := os.Open("scopus.key")
	if err != nil {
		log.Fatal(err)
	}
	apiKey, err := io.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv(paths.ENV_SCOPUS_API_KEY, string(apiKey))
	keyFile.Close()
}

const (
	SERVER_PORT = ":3000"
	QUERY_PORT  = ":3001"
)

func runServices() {
	os.Setenv(paths.ENV_QUERY_PORT, QUERY_PORT)
	setupScopus()

	// TODO prefix logs with process name, so we can tell them apart
	//      maybe use contexts?
	go query.Main()
	go poll.Main()
}

func main() {
	mongod := runMongod()
	defer mongod.Kill()

	go runServices()

	log.Fatal(http.ListenAndServe(SERVER_PORT, http.FileServer(http.Dir("./static"))))
}
