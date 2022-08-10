package files

import (
	"context"
	"errors"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"strings"
)

var ErrInvalidURL = errors.New("invalid URL format")
var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidYAMLFormat = errors.New("invalid YAML format")

type Repo struct {
	Url     string
	Folders []string
	Info    *git.GitRepository
}

func ReadConfig(file string) ([]Repo, error) {
	var repos map[string][]Repo

	data, err := os.ReadFile(file)
	if err != nil {
		return []Repo{}, ErrInvalidFile
	}
	err = yaml.Unmarshal(data, &repos)
	if err != nil {
		return []Repo{}, ErrInvalidYAMLFormat
	}
	// Checking for invalid url formats
	for _, repo := range repos["repos"] {
		_, err := url.ParseRequestURI(repo.Url)
		if err != nil {
			return []Repo{}, ErrInvalidURL
		}
	}
	return repos["repos"], nil
}

// extract from a repo url the Project and the RepositoryId
// 	url := "https://dev.azure.com/schwarzit/schwarzit.stackit-mongodb/_git/stackit-postgres-artifactory-cleanup"
//  RepositoryId := "stackit-postgres-artifactory-cleanup"
//	Project := "schwarzit.stackit-mongodb"
func (repo *Repo) GetInfo(client git.Client, ctx context.Context) {
	var err error
	url, _ := url.ParseRequestURI(repo.Url)
	urlInfo := strings.Split(url.Path[1:], "/")
	args := git.GetRepositoryArgs{
		RepositoryId: &urlInfo[3],
		Project:      &urlInfo[1],
	}
	repo.Info, err = client.GetRepository(ctx, args)
	fmt.Println(err)
	fmt.Printf("%+v\n", repo.Info)
}
