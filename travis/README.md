
## Steps

Note the [Travic CLI](https://github.com/travis-ci/travis.rb) uses 
[Travis API v2](https://docs.travis-ci.com/api) which is phased out by v3.

Get the repo id we'll need for later queries.
```
> travis raw '/repos/ethereum/go-ethereum' --json | jq .repo.id
1697900
```

```
> travis raw '/repos/ethereum/go-ethereum/builds' --json | jq '[.builds[] | del(.config)]' > builds.json
> cat builds.json | jq '.[-1].number'
> travis raw '/repos/ethereum/go-ethereum/builds?after_number=21201' --json | jq '[.builds[] | del(.config)]' | jq -s add builds.json - > tmp.json
> mv tmp.json builds.json
```

Extract job ids of failed PRs
```
> cat pr_builds.json|jq '[.[]|select(.state=="failed") |{"build_id": .id,job_ids}]' > failed_pr_jobs.json
```

Iterate over job ids and get their data into a jsonl file
```
> travis raw '/jobs/651548554' --json | jq '.job | del(.config)' | jq -c . >> failed_jobs.jsonl
```
