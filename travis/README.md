

Download log files
```
> travisor download -repo ethereum/go-ethereum -build.count 250
```

Extract failures from log files stored locally
```
> travisor stats > failures.csv
```

The emitted csv file has this structure

| Occurrences | Test case | Package |
| - | - | - |
| 52 | TestSimulation | github.com/ethereum/go-ethereum/whisper/whisperv6 |
| 31 | TestBroadcastBlock | github.com/ethereum/go-ethereum/eth |
| ... | ... | ... |
