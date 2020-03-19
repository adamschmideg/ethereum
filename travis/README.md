
Each step builds on the previous one(s). Nothing gets overwritten.

## Download build and job data into csv 
```
travisor -do jobs -n 250 -dir .
```
Result: `builds.csv` and `jobs.csv` populated

## Download logs
Downloads all logs referred to in the file `jobs.csv`
```
travisor -do logs -dir .
```
Pre-requisite: `jobs.csv`
Result: `logs/` folder populated with log files

## Find failures
```
travisor -do failures -dir .
```
Pre-requisite: `logs/` folder with log files
Result: `failures.csv`


