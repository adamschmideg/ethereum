package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shuheiktgw/go-travis"
	"strings"
)

func failedBuilds(repo string, n int, buildStates []string, jobStates []string) {
	//c := make(chan travis.Log)
	jobsOpts := travis.JobsOption{State:jobStates}
	client := travis.NewClient(travis.ApiOrgUrl, "")
	offset := 0
	limit := 2
	for offset < n {
		buildOpts := travis.BuildsByRepoOption{Offset:offset, Limit:limit, State: buildStates}
		builds, _, _ := client.Builds.ListByRepoSlug(context.Background(), repo, &buildOpts)
		for _, b := range builds {
			jobs, _, _ := client.Jobs.List(context.Background(), &jobsOpts)
			fmt.Println("build", *b.StartedAt, *b.State, len(jobs))
			for _, j := range jobs {
				log, _, _ := client.Logs.FindByJobId(context.Background(), *j.Id)
				fmt.Println(*log.Href)
			}
		}
		offset += limit
	}
}

func main() {
	repo := flag.String("repo", "ethereum/go-ethereum", "Github <username>/<repo>")
	buildCount := flag.Int("build.count", 2, "Number of builds to check")
	bs := flag.String("build.states", "failed,errored", "Comma-separated list of build states")
	buildStates := strings.Split(*bs, ",")
	js := flag.String("job.states", "failed,errored", "Comma-separated list of job states")
	jobStates := strings.Split(*js, ",")
	flag.Parse()
	failedBuilds(*repo, *buildCount, buildStates, jobStates)
}



