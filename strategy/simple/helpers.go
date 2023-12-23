package simple

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func createFileIfNotExist(file string) (bool, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			return false, err
		}
		log.Printf("created file %s\n", file)
		return false, nil
	}
	return true, nil
}

func createDirIfNotExist(dir string) (bool, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return false, err
		}
		log.Printf("created directory %s\n", dir)
		return false, nil
	}
	return true, nil
}

func dumpState(state *State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("state/%s.json", state.Symbol), data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func percentChange(a, b float64) float64 {
	return ((a - b) / a) * 100
}
