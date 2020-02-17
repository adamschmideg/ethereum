
## Steps

Note the [Travic CLI](https://github.com/travis-ci/travis.rb) uses 
[Travis API v2](https://docs.travis-ci.com/api) which is phased out by v3.

Get the repo id we'll need for later queries.
```
> travis raw '/repos/ethereum/go-ethereum' --json | jq .repo.id
1697900
```

```
> travis raw '/repos/ethereum/go-ethereum/builds?event_type=pull_request' --json | jq .builds > builds.json
> cat builds.json | jq '.[-1].id'
```