package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
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

type travisJob struct {
	build *travis.Build
	job   *travis.Job
}

type travisInfo struct {
	client *travis.Client
	repo   *string
	count  int
	dir *string
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
			builds, _, _ := tr.client.Builds.ListByRepoSlug(context.Background(), *tr.repo, &buildOpts)
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

func (tr *travisInfo) doJobs() error {
	// Create buildsFile
	f := filepath.Join(*tr.dir, "builds.csv")
	if _, err := os.Stat(f); err == nil {
		return fmt.Errorf("File %v exists", f)

	}
	buildsFile, err := os.Create(f)
	if err != nil {
		return err
	}
	defer buildsFile.Close()
	buildsFile.WriteString("Id,Number,State,Duration,EventType,PullRequestNumber,StartedAt,FinishedAt,BranchName,CommitSha,CreatedBy\n")

	// Create jobsFile
	f = filepath.Join(*tr.dir, "jobs.csv")
	if _, err := os.Stat(f); err == nil {
		return fmt.Errorf("File %v exists", f)

	}
	jobsFile, err := os.Create(f)
	if err != nil {
		return err
	}
	defer jobsFile.Close()
	jobsFile.WriteString("Id,Number,BuildId,Number,State,StartedAt,FinishedAt")

	// Process builds and logs
	for job := range tr.buildsAndJobs() {
		if job.build != nil {
			b := job.build
			s := row(b.Id, b.Number, b.State, b.Duration, b.EventType, b.PullRequestNumber,
				b.StartedAt, b.FinishedAt, b.Branch.Name, b.Commit.Sha, b.CreatedBy.Login)
			buildsFile.WriteString(s + "\n")
		}
		if job.job != nil {
			j := job.job
			s := row("job", j.Id, j.Number, j.Build.Id, j.Number, j.State, j.StartedAt, j.FinishedAt)
			jobsFile.WriteString(s + "\n")
		}
	}
	return nil
}

func (tr *travisInfo) doLogs() error {
	return nil
}

func (tr *travisInfo) doFailures() error {
	return nil
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
	return allStats
}

func main() {
	count := flag.Int("n", 3, "Number of builds")
	dir := flag.String("dir", ".", "Working directory")
	action := flag.String("do", "", "Choose one of: jobs|logs|failures")
	repo := flag.String("repo", "ethereum/go-ethereum", "Github <username>/<repo>")
	flag.Parse()

	var err error
	tr := travisInfo{travis.NewClient(travis.ApiOrgUrl, ""), repo, *count, dir}
	switch *action {
	case "jobs":
		err = tr.doJobs()
	case "logs":
		err = tr.doLogs()
	case "failures":
		err = tr.doFailures()
	default:
		err = fmt.Errorf("Uknown action: %v", *action)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func stringify(a interface{}) string {
	raw := spew.Sprint(a)
	raw = strings.Replace(raw, "<*>", "", -1)
	raw = strings.Replace(raw, "<nil>", "", -1)
	return raw
}

func row(a ...interface{}) string {
	var all []string
	for _, arg := range a {
		all = append(all, stringify(arg))
	}
	return strings.Join(all, ",")
}

func newmain() {
	tr := travisInfo{client: travis.NewClient(travis.ApiOrgUrl, "")}
	for job := range tr.buildsAndJobs() {
		if job.build != nil {
			b := job.build
			s := row("build", b.Id, b.Number, b.State, b.Duration, b.EventType, b.PullRequestNumber,
				b.StartedAt, b.FinishedAt, b.Branch.Name, b.Commit.Sha, b.CreatedBy.Login)
			fmt.Println(s)
		}
		if job.job != nil {
			j := job.job
			s := row("job", j.Id, j.Number, j.Build.Id, j.Number, j.State, j.StartedAt, j.FinishedAt)
			fmt.Println(s)

		}
	}
}
