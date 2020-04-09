package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// docker-compose run light
	// rpc call it
	cmd := exec.Command("docker-compose", "run", "lightserver")
	cmd.Env = append(os.Environ(), fmt.Sprintf("RPCPORT=%v", 8548))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
