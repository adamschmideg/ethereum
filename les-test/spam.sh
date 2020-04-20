#!/bin/sh
alias send='curl -s -H "Content-Type: application/json" -X POST --data '
enode=`send '{"jsonrpc":"2.0","id":555,"method":"admin_nodeInfo","params":[]}' http://lightserver:8444 | jq -r .result.enode`
json="{\"jsonrpc\":\"2.0\",\"id\":555,\"method\":\"admin_addPeer\",\"params\":[\"${enode}\"]}"
echo "sending: $json"
send $json http://client_1:9444
send $json http://client_2:9444

ethspam | while read q
do
	# echo $q
	send $q http://client_1:9444
	send $q http://client_2:9444
done
