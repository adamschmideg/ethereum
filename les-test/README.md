```
> light servers -n 3 -networkid 1
# Time passes
> light clients -n 5 -servers 3 -spam -networkid 1
```

Development
```
RPCPORT=8548 NETWORKID=5 docker-compose -f server.yml run --service-ports --no-deps lightserver
RPCPORT=9548 NETWORKID=5 docker-compose -f client.yml run --service-ports client
# In another terminal
curl http://localhost:8548 -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","id":7151823,"method":"web3_clientVersion","params":[]}'
```