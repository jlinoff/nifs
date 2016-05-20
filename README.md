# nifs
Tool to print out a simple report of the IP addresses associated with network interfaces written in Go.

It is somewhat platform independent in the sense that it works on my Mac and various linux distros. I have not tried it on windows because I do not have access to a windows machine.

I wrote it to make parsing interface/IP address simpler for specific tasks that I am working on. It is much easier to parse than ifconfig or "ip addr show" but it doesn't contain nearly as much useful information.

As mentioned above, the report format is very simple. There is one line per interface/IP address combination. Each line has fields described in the table below.

| Field | Description |
| ----- | ----------- |
| interface name | The interface name. Example: lo0 |
| IP address or CIDR | The IP address. Example: 127.0.0.1 |
| unicast or multicast | The broadcast type. Disabled by -c. |
| MAC address | The hardware address. Enabled by -M. |

These are the available options.

| Option | Long | Description |
| ------ | ---- | ----------- |
| -4 | --ipv4 | Only report IPv4 addresses. The default is to report IPv4 and IPv6 addresses. |
| -6 | --ipv6 | Only report IPv6 addresses. The default is to report IPv4 and IPv6 addresses. |
| -h | --help | Print the help message. |
| -H | --hw | Only report HW interfaces (with MAC addresses). The default is to report HW and SW interfaces. |
| -m | --multicast | Only report multicast addresses. The default is to report the unicast and multicast. |
| -S | --SW | Only report SW interfaces (no MAC addresses). The default is to report HW and SW interfaces. |
| -u | --unicast | Only report unicast addresses. The default is to report the unicast and multicast. |
| -V | --version | Print the program version and exit. |


Here is an example run:

```bash
$ # Only display unicast, IPv4 addresses.
$ go run nifs.go -4 -u
lo0 127.0.0.1/8 unicast NOMAC
en0 192.168.200.202/24 unicast a4:5e:60:ef:6a:75
utun0 10.213.42.74/24 unicast NOMAC
```
