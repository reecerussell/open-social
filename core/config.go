package core

import (
	"encoding/json"
	"io/ioutil"
)

// ReadConfig unmarshals JSON from a file to config.
func ReadConfig(filename string, config interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	return nil
}
