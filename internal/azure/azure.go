package azure

import (
	"context"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"kustocheck/internal/files"
	"log"
	"strings"
)

func GetYAMLUrls(repo *files.Repo, gitClient git.Client, ctx context.Context) {
	blob := setBlob(*repo)
	gitItems, err := gitClient.GetItems(ctx, blob)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range *(gitItems) {
		for j := 0; j < len(repo.Folders); j++ {
			if strings.HasPrefix(*(v.Path), repo.Folders[j]) {
				if strings.HasSuffix(*(v.Path), ".yaml") || strings.HasSuffix(*(v.Path),
					".yml") {
					filename := *(repo.Info.Name) + *(v.Path)
					file := files.File{
						Path: "downloads/" + filename,
						Url:  *(v.Url),
					}
					repo.Files = append(repo.Files, file)
				}
			}
		}
	}
	//errch := make(chan error, len(repo.Files))
	//defer close(errch)
	//var wg sync.WaitGroup
	//
	//for _, file := range repo.Files {
	//	wg.Add(1)
	//	go func(file files.File) {
	//
	//		errch <- file.Download()
	//		wg.Done()
	//	}(file)
	//}
	//wg.Wait()
}

func setBlob(repo files.Repo) git.GetItemsArgs {
	repoId := fmt.Sprintf("%v", repo.Info.Id)
	isTrue := true
	isFull := git.VersionControlRecursionTypeValues.Full
	blob := git.GetItemsArgs{
		RepositoryId:   &repoId,
		RecursionLevel: &isFull,
		Download:       &isTrue,
	}

	return blob
}
