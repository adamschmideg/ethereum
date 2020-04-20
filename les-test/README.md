### The `addpeer` command

```
addpeer <node port> <peer port>
```


```
light spam --clients 3 --networkid 5`
```
It will
- Start a server with `docker run`
- Start 3 clients with `docker run` at consecutive ports starting from 9548
- For each client add the server as peer
- Send random JSON-RPC requests to the clients 

### Development
```
 NETWORKID=5 docker-compose up --scale influxdb=0 --scale grafana=0 --build
# In another terminal
curl http://localhost:9444 -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","id":7151823,"method":"web3_clientVersion","params":[]}'
```
To add more clients

 - Add client_n to `docker-compose.yml`
 - Add client_n to `spam.sh` 