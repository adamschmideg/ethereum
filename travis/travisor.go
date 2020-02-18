package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

const (
	baseUrl = "https://api.travis-ci.org"
)

type buildInfo []interface{}
//type buildInfo map[string]interface{}

func readJSONFromUrl(path string) ([]byte, error) {
	url := baseUrl + path
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func builds(repo string, maxBuilds int) (buildInfo,error) {
	var jsonBytes []byte
	var err error
	path := fmt.Sprintf("/repos/%v/builds", repo)
	if jsonBytes, err = readJSONFromUrl(path); err != nil {
		return nil, err
	}
	var b buildInfo
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		return nil, err
	}
	return b, nil
}

func main() {
	repo := flag.String("repo", "ethereum/go-ethereum", "<username>/<repo> on github")
	maxBuilds := flag.Int("max-builds", 5, "Max number of builds")
	flag.Parse()
	b, err := builds(*repo, *maxBuilds)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("builds: %#v", b)
}

