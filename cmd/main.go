package main

import (
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"kustocheck/internal/azure"
	"kustocheck/internal/files"
	"log"
	"os"
)

var configFile = "config/repos.yaml"

const (
	organizationUrl = "https://dev.azure.com/schwarzit"
)

func main() {
	var dependency files.Dependency

	repos, err := files.ReadConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}
	pat := os.Getenv("ADO_TOKEN")
	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(organizationUrl, pat)

	ctx := context.Background()

	// Create a new client to interact with git area
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	files.InitRepo(repos, gitClient, ctx)
	for i := 0; i < len(repos); i++ {
		azure.GetYAMLUrls(&repos[i], gitClient, ctx)
	}
	files.Download(repos)
	dependency.Add("downloads")

}
