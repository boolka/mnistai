package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/boolka/goai/pkg/network"
	"github.com/boolka/mnistai/pkg/config"
	"github.com/boolka/mnistdb/pkg/mnistdb"
	"github.com/boolka/mnistidx/pkg/mnistidx"
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
		pathErr := &os.PathError{}

		if errors.As(err, &pathErr) {
			fmt.Println("network created")
			network.Randomize(netCfg.RandomRate)
		} else {
			log.Fatal(err)
		}
	}

	if netLoadFile != nil {
		defer netLoadFile.Close()

		if err := network.Deserialize(netLoadFile); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("network loaded")
		}
	}

	fmt.Println("network layers:", netCfg.Layers, "epochs:", netCfg.Epochs, "learning rate:", netCfg.LearningRate, "random rate:", netCfg.RandomRate, "skip rate:", netCfg.SkipRate)
	fmt.Println("training...")

	startTraining := time.Now()

	for i := range netCfg.Epochs {
		fmt.Println("epoch:", i)

		trainIdx, err := mnistidx.NewIDX(bytes.NewReader(mnistdb.TrainImages), bytes.NewReader(mnistdb.TrainLabels))
		if err != nil {
			log.Fatal(err)
		}

		buf := trainIdx.NewBuffer()
		expectedOutputs := make([]float64, 10)

		for {
			if rand.Float64() < netCfg.SkipRate {
				continue // Skip some samples to speed up training
			}

			label, err := trainIdx.Read(buf)
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

			clear(expectedOutputs)
			expectedOutputs[label] = 1.0

			network.Correct(inputs, expectedOutputs, netCfg.LearningRate)
		}
	}

	fmt.Println("train duration:", time.Since(startTraining).String())

	netSaveFile, err := os.Create(netFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer netSaveFile.Close()

	if err := network.Serialize(netSaveFile); err != nil {
		log.Fatal(err)
	}
}
