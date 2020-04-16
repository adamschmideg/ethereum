#!/bin/sh
ethspam | while read q                                                       
do
	echo $q
	curl http://client_1:9444  -H "Content-Type: application/json" -X POST --data $q
	curl http://client_2:9444  -H "Content-Type: application/json" -X POST --data $q
done
