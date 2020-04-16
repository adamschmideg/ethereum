#!/bin/sh
ethspam | while read q                                                       
do
	echo curl http://localhost:9441  -H "Content-Type: application/json" -X POST --data $q
	curl http://localhost:9441  -H "Content-Type: application/json" -X POST --data $q
	curl http://localhost:9442  -H "Content-Type: application/json" -X POST --data $q
done
