
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