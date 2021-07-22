package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func GetManifestFromBytes(data []byte) (*Manifest, error) {
	m := new(Manifest)
	err := yaml.Unmarshal(data, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetManifestFromFile returns manifest struct from file path
func GetManifestFromFile(filePath string) (*Manifest, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return GetManifestFromBytes(data)
}
