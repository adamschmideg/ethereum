#!/bin/sh
alias send='curl -s -H "Content-Type: application/json" -X POST --data '
enode=`send '{"jsonrpc":"2.0","id":555,"method":"admin_nodeInfo","params":[]}' http://lightserver:8444 | jq -r .result.enode`
json="{\"jsonrpc\":\"2.0\",\"id\":555,\"method\":\"admin_addPeer\",\"params\":[\"${enode}\"]}"
echo "sending: $json"
send $json http://client_1:9444
send $json http://client_2:9444
send $json http://client_3:9444
send $json http://client_4:9444
send $json http://client_5:9444
send $json http://client_6:9444
send $json http://client_7:9444
send $json http://client_8:9444
send $json http://client_9:9444
send $json http://client_10:9444

ethspam -r 2 | while read q
do
	# echo $q
	send $q http://client_1:9444
	send $q http://client_2:9444
	send $q http://client_3:9444
	send $q http://client_4:9444
	send $q http://client_5:9444
	send $q http://client_6:9444
	send $q http://client_7:9444
	send $q http://client_8:9444
	send $q http://client_9:9444
	send $q http://client_10:9444
done
