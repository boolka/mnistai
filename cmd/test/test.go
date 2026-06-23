package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/boolka/goai/pkg/network"
	"github.com/boolka/mnistai/pkg/config"
	"github.com/boolka/mnistdb/pkg/mnistdb"
	"github.com/boolka/mnistidx/pkg/mnistidx"
)

const surenessThreshold = 0.9

func main() {
	var netFilename, cfgFilename string
	flag.StringVar(&cfgFilename, "config", "config.json", "ai network config filename")
	flag.StringVar(&cfgFilename, "c", "config.json", "ai network config filename")
	flag.StringVar(&netFilename, "net", "ai.net", "ai network filename to load from")
	flag.StringVar(&netFilename, "n", "ai.net", "ai network filename to load from")

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

	f, err := os.Open(netFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := network.Deserialize(f); err != nil {
		log.Fatal(err)
	}

	fmt.Println("network layers:", netCfg.Layers)
	fmt.Println("testing...")

	testIdx, err := mnistidx.NewIDX(bytes.NewReader(mnistdb.TestImages), bytes.NewReader(mnistdb.TestLabels))
	if err != nil {
		log.Fatal(err)
	}

	buf := testIdx.NewBuffer()

	totalCorrect := 0
	totalSureness := 0
	totalIncorrect := 0

	for {
		label, err := testIdx.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		inputs := make([]float64, len(buf))
		for i := 0; i < len(buf); i++ {
			inputs[i] = float64(buf[i]) / 255.0
		}

		outputs := network.Activate(inputs)

		predictedLabel := 0
		maxOutput := outputs[0]
		for i := 1; i < len(outputs); i++ {
			if outputs[i] > maxOutput {
				maxOutput = outputs[i]
				predictedLabel = i
			}
		}

		if predictedLabel == int(label) {
			totalCorrect++

			if maxOutput > surenessThreshold {
				totalSureness++
			}
		} else {
			totalIncorrect++
		}
	}

	fmt.Println("total correct:", totalCorrect)
	fmt.Println("total sureness:", totalSureness)
	fmt.Println("total incorrect:", totalIncorrect)
	fmt.Println("accuracy:", float64(totalCorrect)/float64(totalCorrect+totalIncorrect)*100, "%")
}
