package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type Semver struct {
	Major, Minor, Patch int
}

var repo = "lector"
var org = "drtbz"
var version Semver

//var uri = "https://api.github.com/repos/"+org+"/"+repo+"/releases"

func main() {
	ctx, client := GitHubSetup()
	release, _, err := LatestRelease(client, ctx, org, repo)
	if err != nil {
		log.Panicf("something went wrong grabbing release: \n %v", err)
	}
	tag := removeOuterQuotes(strings.Replace(strings.ToLower(github.Stringify(release.TagName)), "v", ``, 1))
	re, err := regexp.Compile(`^[\d\.]+$`)
	if err != nil {
		log.Panicf("Something went wrong with regex: \n %v", err)
	}

	// Split and convert the tag
	if match := re.MatchString(tag); match {
		log.Printf("match: %v,  %v", match, tag)
		split := strings.Split(tag, ".")
		for k, v := range split {
			switch k {
			case 0:
				version.Major, err = strconv.Atoi(v)
			case 1:
				version.Minor, err = strconv.Atoi(v)
			case 2:
				version.Patch, err = strconv.Atoi(v)
			}
			if err != nil {
				log.Fatal("error converting version")
			}
		}
		log.Printf("%x.%x.%x", version.Major, version.Minor, version.Patch)
	} else {
		log.Fatal("match: %v,  %v", match, tag)
	}
}

func GitHubSetup() (context.Context, *github.Client) {
	token, tkex := os.LookupEnv("GHTOKEN")

	if !tkex {
		log.Fatal("Couldn't get token from ENV")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client
}

func LatestRelease(client *github.Client, ctx context.Context, owner, repo string) (release *github.RepositoryRelease, resp *github.Response, err error) {
	release, resp, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
	return
}

func removeOuterQuotes(s string) string {
	return regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(s, `$1`)
}