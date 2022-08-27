package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"github.com/drtbz/release-semver/version"
)

func main() {
	var repo = os.Args[1]
	var org = os.Args[2]
	var pullID = 1

	token, tkex := os.LookupEnv("GHTOKEN")
	if !tkex {
		log.Fatal("Couldn't get token from ENV")
	}
	ctx, client := GitHubPATSetup(token)

	release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
	if err != nil {
		log.Panicf("something went wrong grabbing release: \n %v", err)
	}

	ver, _, err := SplitAndConvert(release.TagName)
	if err != nil {
		log.Fatalf("Getting PR failed: \n %v", err)
	}
	// Grab the pull request
	pr, _, err := client.PullRequests.Get(ctx, org, repo, pullID)
	if err != nil {
		log.Fatalf("Getting PR failed: \n %v", err)
	}

	newTagName, err := TagFromPRTitle(github.Stringify(pr.Title), ver)
	if err != nil {
		log.Fatalf("Getting new tag failed: \n %v", err)
	}

	commitish := "main"
	log.Printf("TargetCommitish: %v", commitish)
	newReleaseOpts := &github.RepositoryRelease{
		TagName:         &newTagName,
		Name:            &newTagName,
		TargetCommitish: &commitish,
		//Body: newTag.Message,
	}
	newRelease, resp, err := client.Repositories.CreateRelease(ctx, org, repo, newReleaseOpts)

	if newRelease == nil {
		log.Printf("response: %v", resp.Body)
	}
	if err != nil {
		log.Fatalf("Creating release failed \n %v", err)
	}

}

// Sets up a GitHub Client using a PAT
func GitHubPATSetup(pat string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client
}

func removeOuterQuotes(s string) string {
	return regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(s, `$1`)
}

func SplitAndConvert(tagName *string) (ver version.Semver, s string, err error) {
	t := removeOuterQuotes(strings.Replace(strings.ToLower(*tagName), "v", ``, 1))
	re, err := regexp.Compile(`^[\d\.]+$`)
	if err != nil {
		log.Panicf("Something went wrong with regex: \n %v", err)
	}
	// Split and convert the tag as long as it matches the format v#.#.#
	ver = &version.Version{Major: 0, Minor: 0, Patch: 0}
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

func TagFromPRTitle(n string, v version.Semver) (s string, err error) {
	prTitle := removeOuterQuotes(strings.ToLower(n))
	if match, err := regexp.MatchString(`^#major`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Major PR Tag")
		}
		log.Printf("Matched Major PR Tag on commit msg: %v", prTitle)
		v.IncrementMajor()
	}
	if match, err := regexp.MatchString(`^#minor`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Minor PR Tag")
		}
		log.Printf("Matched Minor PR Tag on commit msg: %v", prTitle)
		v.IncrementMinor()
	}
	if match, err := regexp.MatchString(`^#patch`, prTitle); match {
		if err != nil {
			log.Println("Couldnt Match on Patch PR Tag")
		}
		log.Printf("Matched Patch PR Tag on commit msg: %v", prTitle)
		v.IncrementPatch()
	}
	s = v.ToString()
	log.Printf("New Tag will be %v", s)
	return
}
