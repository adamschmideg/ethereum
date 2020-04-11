package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

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
		log.Fatal("Cannot connect server", err)
	}
	client, err := connect(clientPort)
	if err != nil {
		log.Fatal("Cannot connect client", err)
	}

	nodeInfo, err := server.nodeInfo()
	if err != nil {
		return err
	}
	enode := nodeInfo.Enode
	peerCh := make(chan *p2p.PeerEvent)
	sub, err := server.c.Subscribe(context.Background(), "admin", peerCh, "peerEvents")
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	err = client.c.Call(nil, "admin_addPeer", enode)
	if err != nil {
		return err
	}
	dur := 14 * time.Second
	timeout := time.After(dur)
	select {
	case ev := <-peerCh:
		fmt.Printf("At port %v received event: type=%v, peer=%v", serverPort, ev.Type, ev.Peer)
	case err := <-sub.Err():
		return err
	case <-timeout:
		return fmt.Errorf("Timeout after %v", dur)
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
