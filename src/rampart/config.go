// Package rampart is used to convert primal scheme output to rampart configs
package rampart

import (
	"encoding/json"
	"io/ioutil"
)

// Config is
type Config struct {
	References Reference `json:"reference"`
}

// NewConfig is the constructor
func NewConfig(ref *Reference) *Config {
	return &Config{
		References: *ref,
	}
}

// WriteConfig is a method to write the config to a JSON file on disk
func (Config *Config) WriteConfig(filepath string) error {
	enc, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, enc, 0644)
}

// Reference is
type Reference struct {
	Label     string            `json:"label"`
	Accession string            `json:"accession"`
	Length    int               `json:"length"`
	Genes     map[string]Gene   `json:"genes"`
	Amplicons [][]int           `json:"amplicons"`
	Sequence string             `json:"sequence"`
}

// NewReference is the constructor
func NewReference() *Reference {
	return &Reference{
		Label:     "",
		Accession: "",
		Length:    0,
        Genes:     make(map[string]Gene),
        Amplicons: [][]int{},
		Sequence:  "",
	}
}

// Gene is the minimal info for each gene in the reference
type Gene struct {
	Start  int `json:"start"`
	End    int `json:"end"`
	Strand int `json:"strand"`
}
