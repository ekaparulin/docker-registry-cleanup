package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ekaparulin/docker-registry-cleanup/args"
	"github.com/ekaparulin/docker-registry-cleanup/registry"
)

type ImageTag struct {
	Tag     string
	Digest  string
	Created string
}

func main() {

	c, err := args.NewConfig()
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	r := registry.Registry{Url: *c.DockerRegistry}

	if *c.ListSchemaV1 {
		if *c.Verbose {
			println("Getting all tags...")
		}
		tags, err := r.GetTags(*c.DockerRepository)

		if err != nil {
			print(err)
			os.Exit(1)
		}

		if *c.Verbose {
			println("Getting schema 1 tags...")
		}
		images, err := getSchemaV1Images(r, *c.DockerRepository, tags)
		if err != nil {
			print(err)
			os.Exit(1)
		}
		for _, image := range images {
			res := strings.TrimPrefix(image, "https://")
			println(res)
		}

		return
	}

	today := time.Now()
	expireAfter := today.Add(-24 * 365 * 2 * time.Hour)

	if *c.Verbose {
		println("Getting all tags...")
	}
	tags, err := r.GetTags(*c.DockerRepository)

	if err != nil {
		print(err)
		os.Exit(1)
	}

	if *c.Verbose {
		println("Getting expired tags...")
	}
	expiredTags, err := getExpiredTags(r, *c.DockerRepository, tags, expireAfter, *c.Verbose)
	if err != nil {
		print(err)
		os.Exit(1)
	}

	if *c.Verbose {
		println("Deleting expired tags...")
	}
	for _, tag := range expiredTags {
		fmt.Printf("Deleting: " + tag.Tag + " (" + tag.Created + ") ")
		if *c.DryRun {
			continue
		}

		status, err := r.DeleteImage(*c.DockerRepository, tag.Digest)
		fmt.Println(status)

		if err != nil {
			print(err)
		}

	}

}

func getSchemaV1Images(r registry.Registry, repo string, tags registry.TagsResponse) ([]string, error) {
	var ret []string

	for _, tag := range tags.Tags {
		mf, err := r.GetManifest(repo, tag)
		if err != nil {
			return ret, err
		}
		if mf.SchemaVersion == 1 {
			ret = append(ret, r.Url+"/"+repo+":"+tag)
		}
	}
	return ret, nil
}

func getExpiredTags(r registry.Registry, repo string, tags registry.TagsResponse, expireAfter time.Time, verbose bool) ([]ImageTag, error) {
	var ret []ImageTag

	for _, tag := range tags.Tags {

		manifest, err := r.GetManifest(repo, tag)

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error getting manifest for tag: %s, Error: %v\n", tag, err))
		}

		blob, err := r.GetBlob(repo, manifest.Config.Digest)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error getting blob for manifest: %s, Error: %v\n", manifest.Config.Digest, err))
		}

		isExpired, err := expired(blob.Created, expireAfter)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error validation expireation for create date: %s, Error: %v\n", blob.Created, err))
		}

		if verbose {
			println("Tag: " + tag + " Digest: " + manifest.Config.Digest + " Created: " + blob.Created + " Expired: " + fmt.Sprintf("%v", isExpired))
		}

		if !isExpired {
			continue
		}

		ret = append(ret, ImageTag{Tag: tag, Digest: manifest.Config.Digest, Created: blob.Created})

	}
	return ret, nil

}

func expired(timeStamp string, expireAfter time.Time) (bool, error) {

	var ret bool
	date, err := time.Parse(time.RFC3339, timeStamp)

	if err != nil {
		return ret, err
	}

	ret = date.Before(expireAfter)
	return ret, nil
}
