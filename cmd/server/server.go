package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boolka/goai/pkg/network"
	"github.com/boolka/mnistai/pkg/config"
	"github.com/boolka/mnistai/pkg/server"
)

func main() {

	var netFilename, cfgFilename string
	flag.StringVar(&cfgFilename, "config", "config.json", "ai network config filename")
	flag.StringVar(&cfgFilename, "c", "config.json", "ai network config filename")
	flag.StringVar(&netFilename, "net", "ai.net", "ai network filename to save to")
	flag.StringVar(&netFilename, "n", "ai.net", "ai network filename to save to")

	flag.Parse()

	netCfg := &config.NetworkConfig{}

	cfgFile, err := os.Open(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(cfgFile).Decode(netCfg); err != nil {
		log.Fatal(err)
	}

	network := network.NewNetwork(netCfg.Layers)

	netLoadFile, err := os.Open(netFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer netLoadFile.Close()

	if err := network.Deserialize(netLoadFile); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("network loaded")
	}

	server.NewServer(8080, network)
}
