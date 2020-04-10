package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func startServer(rpcport int, networkid int) error {
	cmd := exec.Command("docker-compose", "-f", "server.yml", "run", "--service-ports", "--no-deps", "lightserver")
	cmd.Env = append(os.Environ(), fmt.Sprintf("RPCPORT=%v", rpcport), fmt.Sprintf("NETWORKID=%v", networkid))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", stdoutStderr)
	return nil
}

func main() {
	port := 8548
	networkid := 5 // goerli
	err := startServer(port, networkid)
	if err != nil {
		log.Fatal(err)
	}
	// rpc call it
}
