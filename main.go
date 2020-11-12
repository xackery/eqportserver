package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	//Version is build number
	Version string
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	fmt.Println("Starting EQPortServer", Version)

	for {
		fmt.Println("starting connection")
		listen()
	}
}

func listen() error {
	data := make([]byte, 512)
	conn, err := net.ListenPacket("udp", ":10000")
	if err != nil {
		return fmt.Errorf("listen loginserver on port 10000: %w", err)
	}
	defer conn.Close()

	_, addr, err := conn.ReadFrom(data)
	if err != nil {
		return fmt.Errorf("read %s: %w", addr.String(), err)
	}

	fmt.Println("got connection from", addr.String())
	hosts := strings.Split(addr.String(), ":")
	if len(hosts) < 2 {
		fmt.Printf("%s is an invalid addr", addr.String())
		return nil
	}

	host := fmt.Sprintf("%s:7000", hosts[0])
	err = pingHost(host)
	if err != nil {
		fmt.Println("zone failed", err.Error())
	}
	host = fmt.Sprintf("%s:9000", hosts[0])
	err = pingHost(host)
	if err != nil {
		fmt.Println("world failed", err.Error())
	}
	host = fmt.Sprintf("%s:5999", hosts[0])
	err = pingHost(host)
	if err != nil {
		fmt.Println("login failed", err.Error())
	}

	return nil
}

func pingHost(host string) error {
	conn, err := net.Dial("udp", host)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("rawr!"))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
