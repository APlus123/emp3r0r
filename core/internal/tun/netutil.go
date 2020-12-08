package tun

import (
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// IsPortOpen is this TCP port open?
func IsPortOpen(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

// ValidateIP is this IP legit?
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidateIPPort check if the host string looks like IP:Port
func ValidateIPPort(to string) bool {
	fields := strings.Split(to, ":")
	if len(fields) != 2 {
		return false
	}
	host := fields[0]
	if !ValidateIP(host) {
		return false
	}
	_, err := strconv.Atoi(fields[1])
	return err == nil
}

// IsTor is the C2 on Tor?
func IsTor(addr string) bool {
	if !strings.HasPrefix(addr, "http://") &&
		!strings.HasPrefix(addr, "https://") {
		return false
	}
	nopath := strings.Split(addr, "/")[2]
	fields := strings.Split(nopath, ".")
	return fields[len(fields)-1] == "onion"
}

// HasInternetAccess does this machine has internet access, if yes, what's its exposed IP?
func HasInternetAccess() bool {
	resp, err := http.Get("http://www.msftncsi.com/ncsi.txt")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	if string(respData) == "Microsoft NCSI" {
		return true
	}
	return false
}
