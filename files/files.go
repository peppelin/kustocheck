package files

import (
	"errors"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

var ErrInvalidURL = errors.New("invalid URL format")
var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidYAMLFormat = errors.New("invalid YAML format")

type config struct {
	Pats  map[string]string
	Repos []struct {
		Url     string
		Folders []string
	}
}

func GetPats(file string) (config, error) {

	result := config{}

	data, err := os.ReadFile(file)
	if err != nil {
		return config{}, ErrInvalidFile
	}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return config{}, ErrInvalidYAMLFormat
	}
	// Checking for invalid url formats
	for _, repo := range result.Repos {
		_, err := url.ParseRequestURI(repo.Url)
		if err != nil {
			return config{}, ErrInvalidURL
		}
	}
	return result, nil
}
