package files

import "testing"

func TestGetPats(t *testing.T) {
	t.Run("Getting repo url", func(test *testing.T) {
		result, _ := GetPats("../config/config.yaml.test")
		got := result.Repos[1].Url
		want := "http://github.com"
		if got != want {
			test.Errorf("got %q want %s", got, want)
		}
	})
	t.Run("Getting repo folders", func(test *testing.T) {
		result, _ := GetPats("../config/config.yaml.test")
		got := len(result.Repos[1].Folders)
		want := 4
		if got != want {
			test.Errorf("got %d want %d", got, want)
		}
	})
	t.Run("Getting error when parsing incorrect url", func(test *testing.T) {
		_, err := GetPats("../config/wrong_url.yaml")
		got := err
		want := ErrInvalidURL
		if got != want {
			test.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Getting error when parsing invalid format YAML", func(test *testing.T) {
		_, err := GetPats("../config/wrong_format.yaml")
		got := err
		want := ErrInvalidYAMLFormat
		if got != want {
			test.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Getting error when opening an missing file", func(test *testing.T) {
		_, err := GetPats("../config/missing.yaml")
		got := err
		want := ErrInvalidFile
		if got != want {
			test.Errorf("got %v want %v", got, want)
		}
	})
}
