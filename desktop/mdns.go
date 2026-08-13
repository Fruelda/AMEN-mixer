package main

import (
	"fmt"
	"net"

	"github.com/hashicorp/mdns"
)

var amenMDNSServer *mdns.Server

func startMDNS() {
	if amenMDNSServer != nil {
		return
	}

	ip, err := getLANIPv4()
	if err != nil {
		fmt.Println(
			"[MDNS] IP ERROR:",
			err,
		)

		return
	}

	service, err := mdns.NewMDNSService(
		"AMEN Mixer",
		"_http._tcp",
		"local.",
		"amen-mixer.local.",
		5173,
		[]net.IP{
			ip,
		},
		[]string{
			"app=amen-mixer",
			"ui-port=5173",
			"ws-port=8081",
		},
	)

	if err != nil {
		fmt.Println(
			"[MDNS] SERVICE ERROR:",
			err,
		)

		return
	}

	server, err := mdns.NewServer(
		&mdns.Config{
			Zone: service,
		},
	)

	if err != nil {
		fmt.Println(
			"[MDNS] SERVER ERROR:",
			err,
		)

		return
	}

	amenMDNSServer = server

	fmt.Println(
		"[MDNS] =================================",
	)

	fmt.Println(
		"[MDNS] AMEN MIXER DISCOVERY",
	)

	fmt.Printf(
		"[MDNS] LAN IP : %s\n",
		ip.String(),
	)

	fmt.Println(
		"[MDNS] HOST   : amen-mixer.local",
	)

	fmt.Println(
		"[MDNS] UI     : http://amen-mixer.local:5173",
	)

	fmt.Println(
		"[MDNS] WS     : ws://amen-mixer.local:8081/ws",
	)

	fmt.Println(
		"[MDNS] =================================",
	)
}

func stopMDNS() {
	if amenMDNSServer == nil {
		return
	}

	if err := amenMDNSServer.Shutdown(); err != nil {
		fmt.Println(
			"[MDNS] SHUTDOWN ERROR:",
			err,
		)

		return
	}

	amenMDNSServer = nil

	fmt.Println(
		"[MDNS] STOPPED",
	)
}

func getLANIPv4() (net.IP, error) {
	if conn, err := net.Dial(
		"udp",
		"1.1.1.1:80",
	); err == nil {
		defer conn.Close()

		if localAddr, ok :=
			conn.LocalAddr().(*net.UDPAddr); ok {

			if ip := localAddr.IP.To4(); ip != nil &&
				!ip.IsLoopback() {

				return ip, nil
			}
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil,
			fmt.Errorf(
				"get interfaces: %w",
				err,
			)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addresses {
			var ip net.IP

			switch value := address.(type) {
			case *net.IPNet:
				ip = value.IP

			case *net.IPAddr:
				ip = value.IP

			default:
				continue
			}

			ip = ip.To4()

			if ip == nil ||
				ip.IsLoopback() ||
				ip.IsLinkLocalUnicast() {
				continue
			}

			return ip, nil
		}
	}

	return nil,
		fmt.Errorf(
			"LAN IPv4 address not found",
		)
}
