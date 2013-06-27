package main
import (
    "github.com/ant0ine/go-json-rest"
    "net/http"
    "github.com/daniel-garcia/shipyard"
    "flag"
    "os"
    "fmt"
)

var controller *shipyard.ServiceController

func GetHost(w *rest.ResponseWriter, req *rest.Request) {
    host, _ := controller.Hosts[req.PathParam("id")]
    w.WriteJson(&host)
}

func Hosts(w *rest.ResponseWriter, req *rest.Request) {
    w.WriteJson(&controller.Hosts)
}


var httpFlag string
var daemonFlag bool
var datapathFlag string
var masterAddress string

func init() {
    flag.StringVar(&httpFlag, "http", ":4979", "Port to run the http server or contact the http server")
    flag.BoolVar(&daemonFlag, "daemon", false, "Run in daemon mode.")
    flag.StringVar(&masterAddress, "master", "", "Address of master.")
    envDatapath := os.Getenv("SHIPYARD_DATAPATH")
    if len(envDatapath) == 0 {
        envDatapath = "/var/lib/shipyard/main.json"
    }
    flag.StringVar(&datapathFlag, "datapath", envDatapath, "Path to the data file used for persisting database.")
}


func updateThisHost(du) {
    go func() {
        thisHost, _ := shipyard.CurrentContextAsHost()
        controller.Hosts[thisHost.Id] = thisHost
        controller.Save()
   

func runServer() {
    var err error
    fmt.Printf("Running server on port %s.\n", httpFlag)
	controller, err = shipyard.NewServiceController(datapathFlag)
    if err != nil { fmt.Printf("Error is %s", err) }

    thisHost, _ := shipyard.CurrentContextAsHost()
    controller.Hosts[thisHost.Id] = thisHost
    controller.Save()

    fmt.Printf("Controller running: %s\n", controller)
    handler := rest.ResourceHandler{}
    handler.SetRoutes(
        rest.Route{"GET", "/hosts", Hosts},
        rest.Route{"GET", "/hosts/:id", GetHost},
    )
    http.ListenAndServe(httpFlag, &handler)
}

func clientRequest() {
    fmt.Printf("Client request.\n")
}

func main() {
    flag.Parse()
    if daemonFlag {
        runServer()
    } else {
        clientRequest()
    }
}

