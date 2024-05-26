package config

import (
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp"
	yaml "gopkg.in/yaml.v2"
)

func SampleConfig() {

	AIBFarmConfig = Config{
		API_ENDPOINT: "http://localhost:21881",
	}

	d, err := yaml.Marshal(&AIBFarmConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	file, err := os.OpenFile(AIBFARM_CONFIG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	err = enc.Encode(AIBFarmConfig)
	if err != nil {
		log.Fatalf("error encoding: %v", err)
	}

}

func LoadConfig() error {
	// var str string
	var err error

	file, err := os.OpenFile(AIBFARM_CONFIG, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("error creating config file %v!", err)
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	err = dec.Decode(&AIBFarmConfig)
	if err != nil {
		err = fmt.Errorf("error loading config file %v!", err)
	} else {
		pp.Print(AIBFarmConfig)
	}

	return err
}
