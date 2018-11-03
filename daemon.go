package main

import (
	"flag"
	"fmt"
	"github.com/honeyscience/honeydipper/dipper"
	"os"
)

var config Config

// Services : a catalog of running services in this daemon process
var Services = map[string]*Service{}

var log = dipper.GetLogger("honeydipper")

func init() {
	flag.Usage = func() {
		fmt.Printf("%v [ -h ] service1 service2 ...\n", os.Args[0])
		fmt.Printf("    Supported services include engie, receiver.\n")
		fmt.Printf("  Note: REPO environment variable is required to specify the bootstrap config.\n")
	}
}

func initEnv() {
	flag.Parse()
	config = Config{initRepo: RepoInfo{}, services: flag.Args()}

	ok := true
	if config.initRepo.Repo, ok = os.LookupEnv("REPO"); !ok {
		log.Fatal("REPO environment variable is required to bootstrap honey dipper")
	}
	if config.initRepo.Branch, ok = os.LookupEnv("BRANCH"); !ok {
		config.initRepo.Branch = "master"
	}
	if config.initRepo.Path, ok = os.LookupEnv("BOOTSTRAP_PATH"); !ok {
		config.initRepo.Path = "/"
	}
}

func start() {
	services := config.services
	if len(services) == 0 {
		services = []string{"engine", "receiver", "operator"}
	}
	for _, service := range services {
		switch service {
		case "engine":
			startEngine(&config)
		case "receiver":
			startReceiver(&config)
		case "operator":
			startOperator(&config)
		default:
			log.Fatalf("'%v' service is not implemented", service)
		}
	}
}

func main() {
	initEnv()
	config.bootstrap(".")
	start()
	config.watch()
}
