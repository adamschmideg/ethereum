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
	"reflect"
	"regexp"
	"strings"
)

type travisLog struct {
	job *travis.Job
	log *travis.Log
}

type travisJob struct {
	build *travis.Build
	job   *travis.Job
}

type travisInfo struct {
	client *travis.Client
	repo   string
	count  int
}

func (tr *travisInfo) buildsAndJobs() chan *travisJob {
	limit := 50
	c := make(chan *travisJob, limit+1)
	offset := 0
	go func() {
		for offset < tr.count {
			if offset+limit > tr.count {
				limit = tr.count - offset
			}
			buildOpts := travis.BuildsByRepoOption{Offset: offset, Limit: limit}
			builds, _, _ := tr.client.Builds.ListByRepoSlug(context.Background(), tr.repo, &buildOpts)
			for _, b := range builds {
				c <- &travisJob{b, nil}
				jobs, _, _ := tr.client.Jobs.ListByBuild(context.Background(), *b.Id)
				for _, j := range jobs {
					c <- &travisJob{nil, j}
				}
			}
			offset += limit
		}
		close(c)
	}()
	return c
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
		path := filepath.Join(*logfolder, file.Name())
		data, err := ioutil.ReadFile(path)
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
	repo := downloadCmd.String("repo", "ethereum/go-ethereum", "Github <username>/<repo>")
	buildCount := downloadCmd.Int("build.count", 3, "Number of builds to check")
	logfolderSave := downloadCmd.String("save.folder", "logs", "Logs are saved to this folder")
	bs := downloadCmd.String("build.states", "failed,errored", "Comma-separated list of build states")
	buildStates := strings.Split(*bs, ",")
	js := downloadCmd.String("job.states", "failed,errored", "Comma-separated list of job states")
	jobStates := strings.Split(*js, ",")

	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	logfolderRead := statsCmd.String("read.folder", "logs", "Logs are read from this folder")

	switch os.Args[1] {
	case "download":
		downloadCmd.Parse(os.Args[2:])
		if err := os.MkdirAll(*logfolderSave, 0744); err != nil {
			fmt.Println("log folder", err)
			return
		}
		for log := range filterLogs(*repo, *buildCount, buildStates, jobStates) {
			path := filepath.Join(*logfolderSave, fmt.Sprintf("%v.log", *log.log.Id))
			fmt.Println("Writing to", path)
			if err := ioutil.WriteFile(path, []byte(*log.log.Content), 0644); err != nil {
				fmt.Println("oops", err)
				return
			}
		}
	case "stats":
		statsCmd.Parse(os.Args[2:])
		for pkg, tests := range statsForLogs(logfolderRead) {
			for test, count := range tests {
				fmt.Printf("%v,%v,%v\n", count, test, pkg)
			}
		}
	}
}

func stringify(a ...interface{}) []string {
	var s []string
	for _, arg := range a {
		var elem string
		switch {
		case arg == nil || reflect.ValueOf(arg).IsNil():
			elem = ""
		case reflect.ValueOf(arg).Kind() == reflect.Ptr:
			elem = fmt.Sprintf("%v", reflect.ValueOf(arg).Elem().Interface())
		default:
			elem = fmt.Sprintf("%v", arg)
		}
		s = append(s, elem)
	}
	return s
}

func row(a ...interface{}) string {
	return strings.Join(stringify(a...), ",")
}

func main() {
	tr := travisInfo{client: travis.NewClient(travis.ApiOrgUrl, ""), repo: "ethereum/go-ethereum", count: 2}
	for job := range tr.buildsAndJobs() {
		if job.build != nil {
			b := job.build
			s := row(b.Id, b.Number, b.State, b.Duration, b.EventType, b.PullRequestNumber,
				b.StartedAt, b.FinishedAt, b.Branch.Name, b.Commit.Sha, b.CreatedBy.Login)
			fmt.Println("build:", s)
		}
		if job.job != nil {
			j := job.job
			s := row(j.Id, j.Number, j.Build.Id, j.Number, j.State, j.StartedAt, j.FinishedAt)
			fmt.Println("job:", s)

		}
	}
}
