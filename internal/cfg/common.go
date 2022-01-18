package cfg

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// FromFile loads the cli Configuration at a given path, returns and error if the file does not exists
// or is malformed
func FromFile(path string) (*Configuration, error) {
	cfgfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Configuration{}
	if err := yaml.Unmarshal(cfgfile, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ToFile(filePath string, cfg *Configuration) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	d := path.Dir(filePath)
	if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(filePath, b, 0644); err != nil {
		return err
	}

	return nil
}
