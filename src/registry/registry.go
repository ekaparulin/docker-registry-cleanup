package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Registry struct {
	Url string
}

type TagsResponse struct {
	Name string
	Tags []string
}

type LayerConfig struct {
	MediaType string
	Size      int
	Digest    string
}

type Manifest struct {
	SchemaVersion int
	MediaType     string
	Config        LayerConfig
	Layers        []LayerConfig
}

type Blob struct {
	Created string
}

func (r *Registry) GetTags(repo string) (TagsResponse, error) {

	var responseStructure TagsResponse

	response, err := http.Get(fmt.Sprintf("%s/v2/%s/tags/list", r.Url, repo))

	if err != nil {
		return responseStructure, err
	}

	responseBinary, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseStructure, err
	}

	err = json.Unmarshal(responseBinary, &responseStructure)

	return responseStructure, err
}

func (r *Registry) GetManifest(repo string, tag string) (Manifest, error) {
	var ret Manifest

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/manifests/%s", r.Url, repo, tag), nil)
	if err != nil {
		return ret, err
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	response, err := client.Do(req)

	if err != nil {
		return ret, err
	}

	responseBinary, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(responseBinary, &ret)

	return ret, err
}

func (r *Registry) GetBlob(repo string, digest string) (Blob, error) {
	var ret Blob

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/blobs/%s", r.Url, repo, digest), nil)
	if err != nil {
		return ret, err
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	response, err := client.Do(req)

	if err != nil {
		return ret, err
	}

	responseBinary, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(responseBinary, &ret)

	return ret, err
}

func (r *Registry) DeleteImage(repo string, digest string) (string, error) {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/manifests/%s", r.Url, repo, digest), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}
