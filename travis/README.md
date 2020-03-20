
Each step builds on the previous one(s). Nothing gets overwritten.

## Download build and job data into csv 
```
travisor -do jobs -n 250 -dir .
```
Depends on: -
Outputs: `builds.csv` and `jobs.csv` 

## Download logs
Downloads all logs referred to in the file `jobs.csv`
```
travisor -do logs -dir .
```
Depends on: `jobs.csv`
Outputs: `logs.csv` and `logs/*.txt`

## Find failures
```
travisor -do failures -dir .
```
Depends on: `logs.csv` and `logs/*.txt`
Outputs: `failures.csv`


