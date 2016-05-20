// List the network interfaces with IP address in a platform independent way.
// License: MIT
// Copyright (c) 2016 by Joe Linoff
package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

// NetworkInterfaceForIP contains the network interface information for a single
// IP address.
type NetworkInterfaceForIP struct {
	Interface net.Interface
	IPAddr    string
	Unicast   bool // true == unicast, false == multicast
	IPv4      bool // true == IPv4, false == IPv4
}

// Options contains the command line options.
type Options struct {
	Verbose       int
	ShowMulticast bool
	ShowUnicast   bool
	ShowIPv4s     bool
	ShowIPv6s     bool
	ShowHW        bool // have MAC
	ShowSW        bool // do not have MAC
}

func main() {
	opts := getOpts()
	nifs := LoadNifs()
	ReportNifs(opts, nifs)
	// report(opts, nifs)
}

// LoadNifs loads the network interface information all IP addresses.
func LoadNifs() []NetworkInterfaceForIP {
	recs := []NetworkInterfaceForIP{}

	nifs, _ := net.Interfaces()
	for _, nif := range nifs {
		// unicast
		ifaddrs, _ := nif.Addrs()
		for _, ifaddr := range ifaddrs {
			addr := ifaddr.String()
			//Can't use net.ParseIP() with To4() because it returns false
			// for 127.0.0.1/8.
			ipv4 := strings.Contains(addr, ".")
			rec := NetworkInterfaceForIP{
				Interface: nif,
				IPAddr:    addr,
				Unicast:   true,
				IPv4:      ipv4,
			}
			recs = append(recs, rec)
		}

		// multicast
		ifaddrs, _ = nif.MulticastAddrs()
		for _, ifaddr := range ifaddrs {
			addr := ifaddr.String()
			ipv4 := strings.Contains(addr, ".")
			rec := NetworkInterfaceForIP{
				Interface: nif,
				IPAddr:    addr,
				Unicast:   false,
				IPv4:      ipv4,
			}
			recs = append(recs, rec)
		}
	}
	return recs
}

// ReportNifs reports the NIF information.
func ReportNifs(opts Options, nifs []NetworkInterfaceForIP) {
	for _, nif := range nifs {
		if opts.ShowIPv4s == false && nif.IPv4 == true {
			continue
		}
		if opts.ShowIPv6s == false && nif.IPv4 == false {
			continue
		}
		if opts.ShowUnicast == false && nif.Unicast == true {
			continue
		}
		if opts.ShowMulticast == false && nif.Unicast == false {
			continue
		}
		if opts.ShowHW == false && len(nif.Interface.HardwareAddr) > 0 {
			continue
		}
		if opts.ShowSW == false && len(nif.Interface.HardwareAddr) == 0 {
			continue
		}

		fmt.Printf("%v ", nif.Interface.Name)
		fmt.Printf("%v ", nif.IPAddr)
		if nif.Unicast {
			fmt.Printf("unicast ")
		} else {
			fmt.Printf("multicast ")
		}
		mac := nif.Interface.HardwareAddr.String()
		if len(mac) == 0 {
			mac = "NOMAC"
		}
		fmt.Printf("%v", mac)
		fmt.Printf("\n")
	}
}

// Get the command line options.
func getOpts() Options {
	opts := Options{
		ShowUnicast:   false,
		ShowMulticast: false,
		ShowIPv4s:     false,
		ShowIPv6s:     false,
		ShowHW:        false,
		ShowSW:        false,
	}
	for i := 1; i < len(os.Args); i++ {
		opt := os.Args[i]
		switch opt {
		case "-h", "--help":
			help()
			os.Exit(0)
		case "-H", "--hw":
			opts.ShowHW = true
		case "-m", "--multicast":
			opts.ShowMulticast = true
		case "-u", "--unicast":
			opts.ShowUnicast = true
		case "-4", "--ipv4":
			opts.ShowIPv4s = true
		case "-6", "--ipv6":
			opts.ShowIPv6s = true
		case "-s", "--sw":
			opts.ShowSW = true
		case "-V", "--version":
			base := path.Base(os.Args[0])
			fmt.Printf("%s 0.2.0\n", base)
			os.Exit(0)
		default:
			fmt.Printf("ERROR: unrecognized option '%v'\n", opt)
			os.Exit(1)
		}
	}
	if opts.ShowIPv4s == false && opts.ShowIPv6s == false {
		opts.ShowIPv4s = true
		opts.ShowIPv6s = true
	}
	if opts.ShowUnicast == false && opts.ShowMulticast == false {
		opts.ShowUnicast = true
		opts.ShowMulticast = true
	}
	if opts.ShowHW == false && opts.ShowSW == false {
		opts.ShowHW = true
		opts.ShowSW = true
	}

	return opts
}

// Get the next argument.
func getNextArg(i *int) string {
	*i++
	if *i >= len(os.Args) {
		fmt.Printf("ERROR: missing argument for %s", os.Args[*i-1])
		os.Exit(1)
	}
	return os.Args[*i]
}

// help.
func help() {
	base := path.Base(os.Args[0])
	msg := `
USAGE
    %[1]v [OPTIONS]

DESCRIPTION
    Generates a very simple report of the IP addresses associated with the
    network interfaces.

    The report format is very simple. There is one line per interface/ip address
    combination. There are 4 fields on each line:

       1. interface name
       2. IP address (or CIDR)
       3. unicast or multicast
       4. MAC address

OPTIONS
    -4, --ipv4         Report only IPv4 addresses.
                       The default is to report both IPv4 and IPv6 addresses.

    -6, --ipv6         Report only IPv6 addresses.
    The default is to report both IPv4 and IPv6 addresses.

    -h, -help          This help message.

    -H, --hw           Report hardware only interfaces (have a MAC addresses).
                       The default is to report all HW and SW interfaces.

    -m, --multicast    Report only multicast addresses.
                       The default is to report unicast and multicast.

    -s, --sw           Report software only interfaces(no MAC addresses).
                       The default is to report all HW and SW interfaces.

    -u, --unicast      Report only unicast addresses.
                       The default is to report unicast and multicast.

    -V, --version      Print the program version and exit.

EXAMPLES
    $ # Example 1: help
    $ %[1]s -h

    $ # Example 2: Report all IP address
    $ %[1]s

    $ # Example 3: Report only IPv4, unicast addresses.
    $ %[1]s -4 -u

    $ # Example 4: Report only HW, IPv4, unicast addresses
    $ %[1]s -4 -u -H

`
	fmt.Printf(msg, base)
}
