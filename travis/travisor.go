package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shuheiktgw/go-travis"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	failingTestRE := regexp.MustCompile("^--- FAIL: (\\w+).*")
	failingPkgRE := regexp.MustCompile("^FAIL\\s+(\\S+)\\s.*")
	var s = make(Stats)
	var names []string
	var pkg string
	for _, line := range strings.Split(*logcontent, "\n") {
		badPkg := failingPkgRE.FindStringSubmatch(line)
		if len(badPkg) > 0 {
			pkg = badPkg[1]
			// Now we know the package, copy all test names into the result struct
			s[pkg] = make(map[string]int)
			for _, n := range names {
				s[pkg][n] = 1
			}
			// Empty names for the next round
			names = []string{}
			continue
		}
		testName := failingTestRE.FindStringSubmatch(line)
		if len(testName) > 0 {
			names = append(names, testName[1])
			continue
		}
	}
	return s
}

// Change `all`
func combineStats(all *Stats, one *Stats) {
	allStats := *all
	oneStats := *one
	for pkg, tests := range oneStats {
		earlierTests, pkgOk := allStats[pkg]
		if pkgOk {
			for testName, _ := range tests {
				earlierTests[testName] += 1
			}
		} else {
			allStats[pkg] = tests
		}
	}
}

func statsForLogs(logfolder *string) Stats {
	allStats := make(Stats)
	files, err := ioutil.ReadDir(*logfolder)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		data, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Fatal(err)
		}
		content := string(data)
		newStats := GetStats(&content)
		combineStats(&allStats, &newStats)
	}
	return  allStats
}

func main() {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	repo := flag.String("repo", "ethereum/go-ethereum", "Github <username>/<repo>")
	buildCount := flag.Int("build.count", 3, "Number of builds to check")
	logfolder := flag.String("log.folder", "logs", "Logs are save to and read from this folder")
	bs := flag.String("build.states", "failed,errored", "Comma-separated list of build states")
	buildStates := strings.Split(*bs, ",")
	js := flag.String("job.states", "failed,errored", "Comma-separated list of job states")
	jobStates := strings.Split(*js, ",")

	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	logfolder = flag.String("log.folder", "logs", "Logs are save to and read from this folder")

	switch os.Args[1] {
	case "download":
		downloadCmd.Parse(os.Args[2:])
		if err := os.MkdirAll(*logfolder, 0744); err != nil {
			fmt.Println("log folder", err)
			return
		}
		for log := range filterLogs(*repo, *buildCount, buildStates, jobStates) {
			path := filepath.Join(*logfolder, fmt.Sprintf("%v.log", *log.log.Id))
			fmt.Println("Writing to", path)
			if err := ioutil.WriteFile(path, []byte(*log.log.Content), 0644); err != nil {
				fmt.Println("oops", err)
				return
			}
		}
	case "stats":
		statsCmd.Parse(os.Args[2:])
		s := statsForLogs(logfolder)
		fmt.Println("stats:", s)
	}
}

