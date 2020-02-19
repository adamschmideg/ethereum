package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shuheiktgw/go-travis"
	"strings"
)

type travisLog struct {
	job *travis.Job
	log *travis.Log
}

func filterLogs(repo string, n int, buildStates []string, jobStates []string) chan travisLog {
	c := make(chan travisLog)
	jobStatesMap := make(map[string]bool)
	for _, k := range jobStates {
		jobStatesMap[k] = true
	}
	client := travis.NewClient(travis.ApiOrgUrl, "")
	offset := 0
	limit := 50
	go func() {
		for offset < n {
			if offset + limit > n {
				limit = n - offset
			}
			buildOpts := travis.BuildsByRepoOption{Offset:offset, Limit:limit, State: buildStates}
			builds, _, _ := client.Builds.ListByRepoSlug(context.Background(), repo, &buildOpts)
			for _, b := range builds {
				jobs, _, _ := client.Jobs.ListByBuild(context.Background(), *b.Id)
				for _, j := range jobs {
					if want, _ := jobStatesMap[*j.State]; want {
						log, _, _ := client.Logs.FindByJobId(context.Background(), *j.Id)
						c <- travisLog{j, log}
					}
				}
			}
			offset += limit
		}
		close(c)
	}()
	return c
}

type Stats map[string]map[string]int

func GetStats(logcontent *string) Stats {
	var s Stats
	return s
}

func testStats(repo string, n int) Stats {
	states := []string{"failed"}
	s := make(Stats)
	for log := range filterLogs(repo, n, states, states) {
		for pkg, tests := range GetStats(log.log.Content) {
			earlierTests, pkgOk := s[pkg]
			if pkgOk {
				for testName, _ := range tests {
					earlierTests[testName] += 1
				}
			}
		}
	}
	return s
}

func main() {
	repo := flag.String("repo", "ethereum/go-ethereum", "Github <username>/<repo>")
	buildCount := flag.Int("build.count", 3, "Number of builds to check")
	bs := flag.String("build.states", "failed,errored", "Comma-separated list of build states")
	buildStates := strings.Split(*bs, ",")
	js := flag.String("job.states", "failed,errored", "Comma-separated list of job states")
	jobStates := strings.Split(*js, ",")
	flag.Parse()
	for log := range filterLogs(*repo, *buildCount, buildStates, jobStates) {
		fmt.Println(*log.job.Build.Number, *log.job.Number, *log.log.Href)
	}
}

