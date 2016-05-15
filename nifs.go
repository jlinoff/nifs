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

func main() {
	opts := getOpts()
	report(opts)
}

// Report the interfaces.
func report(opts Options) {
	nifs, _ := net.Interfaces()
	for _, nif := range nifs {
		if opts.ShowUnicast {
			ifaddrs, _ := nif.Addrs()
			for _, ifaddr := range ifaddrs {
				addr := ifaddr.String()
				printEntry(opts, nif, addr, "unicast")
			}
		}

		if opts.ShowMulticast {
			ifaddrs, _ := nif.MulticastAddrs()
			for _, ifaddr := range ifaddrs {
				addr := ifaddr.String()
				printEntry(opts, nif, addr, "multicast")
			}
		}
	}
}

func printEntry(opts Options, nif net.Interface, addr string, cast string) {
	// Can't use net.ParseIP() with To4() because it returns false
	// for 127.0.0.1/8.
	if opts.ShowIPv4s == false {
		if strings.Contains(addr, ".") {
			return // dumb test.
		}
	}
	if opts.ShowIPv6s == false {
		if strings.Contains(addr, ":") {
			return // dumb test.
		}
	}
	fmt.Printf("%v ", nif.Name)
	fmt.Printf("%v ", addr)

	if opts.PrintCast {
		fmt.Printf("%v ", cast)
	}

	// Print MAC adddresses if they asked for it.
	// If the interface has nbo MAC address print NOMAC to make it easier to
	// recognize.
	if opts.PrintMacs {
		if len(nif.HardwareAddr) > 0 {
			fmt.Printf("%v ", nif.HardwareAddr)
		} else {
			fmt.Printf("%v ", "NOMAC")
		}
	}
	fmt.Printf("\n")
}

// Options contains the command line options.
type Options struct {
	Verbose       int
	ShowMulticast bool
	ShowUnicast   bool
	ShowIPv4s     bool
	ShowIPv6s     bool
	PrintMacs     bool
	PrintCast     bool
}

// Get the command line options.
func getOpts() Options {
	opts := Options{
		ShowUnicast:   true,
		ShowMulticast: true,
		ShowIPv4s:     true,
		ShowIPv6s:     true,
		PrintMacs:     false,
		PrintCast:     true,
	}
	for i := 1; i < len(os.Args); i++ {
		opt := os.Args[i]
		switch opt {
		case "-h", "--help":
			help()
			os.Exit(0)
		case "-c", "--no-cast":
			opts.PrintCast = false
		case "-m", "--no-multicast":
			opts.ShowMulticast = false
		case "-u", "--no-unicast":
			opts.ShowUnicast = false
		case "-4", "--no-ipv4s":
			opts.ShowIPv4s = false
		case "-6", "--no-ipv6s":
			opts.ShowIPv6s = false
		case "-M", "--print-macs":
			opts.PrintMacs = true
		case "-V", "--version":
			base := path.Base(os.Args[0])
			fmt.Printf("%s 1.0.0\n", base)
			os.Exit(1)
		default:
			fmt.Printf("ERROR: unrecognized option '%v'\n", opt)
			os.Exit(1)
		}
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
    combination. There are 2, 3 or 4 fields on each line depdending on whether or
    not you specified -M (--print-macs) or -c (--no-cast).

    These are the fields:

       1. interface name
       2. IP address (or CIDR)
       3. optional unicast or multicast, if -c is specified it won't be printed
       4. optional MAC address if -M was specified, if the interface is not
          associated with a physical device NOMAC is printed.

OPTIONS
    -4, --no-ipv4s     Do not report IPv4 addresses.
                       The default is to report them.

    -6, --no-ipv6s     Do not report IPv6 addresses.
                       The default is to report them.

    -c, --no-cast      Do not print unicast or multicast.
                       The default is to report them.

    -h, -help          This help message.

    -m, --no-multicast Do not report multicast addresses.
                       The default is to report them.

    -M, --print-macs   Print MAC addresses.
                       The default is to not report them.

    -u, --no-unicast   Do not report unicast addresses.
                       The default is to report them.

    -V, --version      Print the program version and exit.

EXAMPLES
    $ # Example 1: help
    $ %[1]s -h

    $ # Example 2: print only the IPv4, unicast interface IP addresses
    $ #            This is done by disabling IPv6 and multicast.
    $ %[1]s -6 -m -c
    lo0 127.0.0.1/8
    en9 192.168.2.104/24
    utun0 172.168.23.87/24

    $ # Example 3: Print all of the interfaces with IP addreses.
    $ %[1]s
    <output snipped>

`
	fmt.Printf(msg, base)
}
