package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"strings"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const(
	GITHUB_TOKEN = "GITHUB_TOKEN"
	SSH_PRIVATE_KEY_PATH = "SSH_PRIVATE_KEY_PATH"
)

var config = viper.New()

func main() {

	config.SetDefault(GITHUB_TOKEN, "placeholder")
	config.SetDefault(SSH_PRIVATE_KEY_PATH, "placeholder")
	config.AutomaticEnv()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GetString(GITHUB_TOKEN)},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 200},
	}

	// list all repositories for the authenticated user
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil{
			logrus.Errorf("ERROR listing repos: %v", err)
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

	}


	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Repository Name", "Branch", "Lines Qty"})

	var counter = 0

	for _, r := range allRepos{

		repoName := *r.Name
		repoGitUrl := convertHttpsToSshUrl(*r.CloneURL)
		repoCloneOutputPath := filepath.FromSlash(`bin\repos\`) + repoName
		defaultBranch := *r.DefaultBranch

		logrus.Infof("INFO Analyzing repository {%s} default branch {%s}", repoName, defaultBranch)

		if cloneRepo(repoName, repoGitUrl, repoCloneOutputPath) {

			var linesTotal = 0
			e := filepath.Walk(repoCloneOutputPath, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() {
					if !strings.Contains(path, string(filepath.Separator)+".git"+string(filepath.Separator)) && !strings.HasSuffix(path, string(filepath.Separator)+".git") {
						//fmt.Println(path)

						count, ok := getLinesCount(path)
						if ok {
							linesTotal += count
						}

						logrus.Infof("INFO analyzing file {%s}: {%d} lines", path, count)
					}
				}
				return err
			})

			logrus.Infof("INFO total lines for repository {%s} branch {%s}: {%d} lines", repoName, defaultBranch, linesTotal)

			if e != nil {
				logrus.Errorf("ERROR walking through files: %v", e)
			}

			t.AppendRow([]interface{}{
				counter,
				repoName,
				defaultBranch,
				linesTotal,
			})

			counter++
		}
	}

	t.Render()
}

func convertHttpsToSshUrl(gitUrl string) (sshUrl string){
	var u, err = url.Parse(gitUrl)
	if err != nil {
		fmt.Errorf("ERROR parsing URL %s: %v", gitUrl, err)
	}

	return fmt.Sprintf("git@github.com:%s", u.Path[1:])
}

func cloneRepo(repoName, repoUrl, outputPath string) (ok bool){
	pem, _ := ioutil.ReadFile(config.GetString(SSH_PRIVATE_KEY_PATH))
	signer, _ := ssh.ParsePrivateKey(pem)
	aith := &ssh2.PublicKeys{User: "git", Signer: signer}

	_, err := git.PlainClone(outputPath, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              aith,
	})

	if err != nil{
		logrus.Errorf("ERROR cloning repository %s: %v", repoName, err)
		return false
	}

	return true
}

func getLinesCount(filePath string) (count int, ok bool) {

	r, err := os.Open(filePath)
	defer r.Close()

	if err != nil {
		logrus.Errorf("ERROR opening file {%s} for read: %v", filePath, err)
		return count, false
	}

	var readSize int

	buf := make([]byte, 1024)

	for {
		readSize, err = r.Read(buf)
		if err != nil {
			break
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], '\n')
			if i == -1 || readSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
	}
	if readSize > 0 && count == 0 || count > 0 {
		count++
	}
	if err == io.EOF {
		return count, true
	}


	return count, false
}