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

func main() {
	var err error
	port := 8548
	networkid := 5 // goerli
	err = startServer(port, networkid)
	if err != nil {
		log.Fatal(err)
	}

	cli := client{}
	url := fmt.Sprintf("http://127.0.0.1:%v", port)
	cli.c, err = rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	nodeInfo := &p2p.NodeInfo{}
	cli.c.Call(&nodeInfo, "admin_nodeInfo")
	fmt.Println("nodeInfo", nodeInfo)
}
