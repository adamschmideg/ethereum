package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

type client struct {
	c      *rpc.Client
}

func startServer(rpcport int, networkid int) error {
	cmd := exec.Command("docker-compose", "-f", "server.yml", "run", "--service-ports", "--no-deps", "lightserver")
	cmd.Env = append(os.Environ(), fmt.Sprintf("RPCPORT=%v", rpcport), fmt.Sprintf("NETWORKID=%v", networkid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	fmt.Printf("running")
	return nil
}

var portToClient = make(map[int]*client)

func connect(port int) (*client, error) {
	cli, _ := portToClient[port]
	if cli != nil {
		// already connected
		return cli, nil
	}
	cli = &client{}
	url := fmt.Sprintf("http://127.0.0.1:%v", port)
	var err error
	cli.c, err = rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	portToClient[port] = cli
	return cli, nil
}

func (cli *client) nodeInfo() (*p2p.NodeInfo, error) {
	nodeInfo := &p2p.NodeInfo{}
	err := cli.c.Call(&nodeInfo, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return nodeInfo, nil
}

func addPeer(serverPort int, clientPort int) error {
	server, err := connect(serverPort)
	if err != nil {
		return err
	}
	client, err := connect(clientPort)
	if err != nil {
		return err
	}

	nodeInfo, err := server.nodeInfo()
	if err != nil {
		return err
	}
	enode := nodeInfo.Enode
	err = client.c.Call(nil, "admin_addPeer", enode)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	serverPort := 8548
	clientPort := 9548

	err := addPeer(serverPort, clientPort)
	if err != nil {
		log.Fatal(err)
	}
}
