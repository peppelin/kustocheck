package files

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidURL        = errors.New("invalid URL format")
	ErrInvalidFile       = errors.New("invalid file")
	ErrInvalidYAMLFormat = errors.New("invalid YAML format")
	ErrDownloadingFiles  = errors.New("error while downloading files")
)

type Repo struct {
	Url     string
	Folders []string
	Files   []File
	Info    *git.GitRepository
}

type File struct {
	Path string
	Url  string
}

type Dependency struct {
	Origin  []string
	Version string
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
func (repo *Repo) getInfo(client git.Client, ctx context.Context) error {
	var err error
	url, _ := url.ParseRequestURI(repo.Url)
	urlInfo := strings.Split(url.Path[1:], "/")
	args := git.GetRepositoryArgs{
		RepositoryId: &urlInfo[3],
		Project:      &urlInfo[1],
	}
	repo.Info, err = client.GetRepository(ctx, args)
	return err
}

func InitRepo(repos []Repo, client git.Client, ctx context.Context) error {
	errch := make(chan error, len(repos))
	defer close(errch)
	var wg sync.WaitGroup

	for i := 0; i < len(repos); i++ {
		wg.Add(1)
		go func(repo *Repo, client git.Client, ctx context.Context) {
			errch <- repo.getInfo(client, ctx)
			wg.Done()
		}(&repos[i], client, ctx)
	}
	wg.Wait()

	for i := 0; i < len(errch); i++ {
		err := <-errch
		if err != nil {
			return err
		}
	}
	return nil
}

func (file File) download() error {

	pat := os.Getenv("ADO_TOKEN")
	request, err := http.NewRequest(http.MethodGet, file.Url, nil)
	request.SetBasicAuth("", pat)
	client := http.Client{Timeout: 5 * time.Second}

	folder := filepath.Dir(file.Path)
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(file.Path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Download(repos []Repo) {
	//errchRepo := make(chan error, len(repos))
	//defer close(errchRepo)
	var wgRepo sync.WaitGroup

	for i := 0; i < len(repos); i++ {
		wgRepo.Add(1)
		errchFile := make(chan error, len(repos[i].Files))
		defer close(errchFile)
		var wgFile sync.WaitGroup

		for _, file := range repos[i].Files {
			wgFile.Add(1)
			go func(file File) {
				errchFile <- file.download()
				wgFile.Done()
			}(file)
		}
		wgFile.Wait()
		wgRepo.Done()
		for i := 0; i < len(errchFile); i++ {
			err := <-errchFile
			if err != nil {
				log.Println(ErrDownloadingFiles)
			}
		}
	}
	wgRepo.Wait()
	//fmt.Println(len(errchRepo))

	//for i := 0; i < len(errchRepo); i++ {
	//	err := <-errchRepo
	//	if err == nil {
	//		log.Fatal(err)
	//	}
	//}

}

func (dependency *Dependency) Add(path string) error {

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if !strings.Contains(path, "openapi") && !strings.Contains(path, ".azuredevops") {
				// check for github content
				dependency.Origin = append(dependency.Origin, path)
				githubReferences(path)

			}
		}
		return nil
	})
	return err
}

func githubReferences(file string) {
	var repo []string
	var repoName string
	var depen Dependency
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file word by word using scanner
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		// do something with a word
		if strings.Contains(scanner.Text(), "github") {
			_, result, _ := strings.Cut(scanner.Text(), ".com/")
			repo = strings.SplitN(result, "/", 4)
			repoName = fmt.Sprintf("%s/%s", repo[0], repo[1])
			fmt.Printf("%+v\n", repoName)
			GithubVersion(repo[0], repo[1], file)
		}
	}

	depen.Origin = repo
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func GithubVersion(owner, repo, file string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	version := "v1.3.8"
	//for _, repo := range repoMap {
	release, response, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		log.Fatalf("repositories.getlatestrelease status code: %v", response.StatusCode)
	}
	fmt.Printf("File: %s\nRepo:%v\tVersion:%s\tGithub Version:%v\n", file, owner+"/"+repo,
		version, *(release).TagName)

	//}
}
