package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// Testing Vars
var repo = "github-release-semver"
var org = "drtbz"
var pullID = 1
var ver Semver = &Version{Major: 0, Minor: 0, Patch: 0}

func main() {
	ctx, client := GitHubSetup()
	release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
	if err != nil {
		log.Panicf("something went wrong grabbing release: \n %v", err)
	}

	_, err = SplitAndConvert(ver, release.TagName)
	if err != nil {
		log.Fatalf("Getting PR failed: \n %v", err)
	}
	// Grab the pull request
	pr, _, err := client.PullRequests.Get(ctx, org, repo, pullID)
	if err != nil {
		log.Fatalf("Getting PR failed: \n %v", err)
	}

	newTag, err := TagFromPRTitle(github.Stringify(pr.Title), ver)
	if err != nil {
		log.Fatalf("Getting new tag failed: \n %v", err)
	}
	
	commitish := github.Stringify(pr.Head.Ref)
	newReleaseOpts := &github.RepositoryRelease{
		TagName: &newTag,
		Name: &newTag,
		TargetCommitish: &commitish,
	}
	newRelease, resp, err := client.Repositories.CreateRelease(ctx, org, repo, newReleaseOpts)

	if newRelease == nil {
		log.Printf("response: %v", resp.Body)
	}
	if err != nil {
		log.Fatalf("Creating release failed \n %v", err)
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

func removeOuterQuotes(s string) string {
	return regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(s, `$1`)
}

func SplitAndConvert(ver Semver, tagName *string) (s string, err error) {
	t := removeOuterQuotes(strings.Replace(strings.ToLower(*tagName), "v", ``, 1))
	re, err := regexp.Compile(`^[\d\.]+$`)
	if err != nil {
		log.Panicf("Something went wrong with regex: \n %v", err)
	}
	// Split and convert the tag as long as it matches the format v#.#.#
	if match := re.MatchString(t); match {
		split := strings.Split(t, ".")
		for k, v := range split {
			switch k {
			case 0:
				ver.SetMajor(v)
			case 1:
				ver.SetMinor(v)
			case 2:
				ver.SetPatch(v)
			}
			if err != nil {
				log.Fatal("error converting version")
			}
		}
		s = ver.ToString()
		log.Printf("Current tag is: %v", s)
	} else {
		log.Fatalf("match: %v,  %v", match, t)
	}
	return
}

func TagFromPRTitle(n string, v Semver) (s string, err error) {
	prTitle := removeOuterQuotes(strings.ToLower(n))
	if match, err := regexp.MatchString(`^#major`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Major PR Tag")
		}
		log.Printf("Matched Major PR Tag on commit msg: %v", prTitle)
		ver.IncrementMajor()
	}
	if match, err := regexp.MatchString(`^#minor`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Minor PR Tag")
		}
		log.Printf("Matched Minor PR Tag on commit msg: %v", prTitle)
		ver.IncrementMinor()
	}
	if match, err := regexp.MatchString(`^#patch`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Patch PR Tag")
		}
		log.Printf("Matched Patch PR Tag on commit msg: %v", prTitle)
		ver.IncrementPatch()
	}
	s = v.ToString()
	log.Printf("New Tag will be %v", s)
	return
}
