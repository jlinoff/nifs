# nifs
Tool to print out a simple report of the IP addresses associated with network interfaces written in Go.

It is somewhat platform independent in the sense that it works on my Mac and various linux distros. I have not tried it on windows because I do not have access to a windows machine.

I wrote to make parsing interface/IP address simpler for specific tasks that I am working on. It is much easier to parse than ifconfig or "ip addr show" but it doesn't contain nearly as much useful information.

As mentioned above, the report format is very simple. There is one line per interface/IP address combination. Each line has 2, 3 or 4 entries depending on the options specified. The output fields are described in the table below.

| Field | Optional | Description |
| ----- | :--------: | ----------- |
| interface name | No | The interface name. Example: lo0 |
| IP address or CIDR | No | The IP address. Example: 127.0.0.1 |
| unicast or multicast | Yes | The broadcast type. Disabled by -c. |
| MAC address | Yes | The hardware address. Enabled by -M. |

These are the available options.

| Option | Long | Description |
| ------ | ---- | ----------- |
| -4 | --no-ipv4s | Do not report IPv4 addresses. The default is to report them. |
| -6 | --no-ipv6s | Do not report IPv6 addresses. The default is to report them. |
| -c | --no-cast | Do not print unicast or multicast. The default is to report them. |
| -h | --help | Print the help message. |
| -m | --no-multicast | Do not report multicast addresses. The default is to report them. |
| -M | --print-macs | Print MAC addresses. The default is to not report them. |
| -u | --no-unicast | Do not report unicast addresses. The default is to report them. |
| -V | --version | Print the program version and exit. |


Here is an example run:

```bash
$ # Only display unicast, IPv4 addresses, do not print unicast or multicast.
$ # This is done by telling it to ignore multicast, ignore IPv6 and disable
$ # the printing of unicast.
$ go run nifs.go -m -6 -c
lo0 127.0.0.1/8
en9 192.168.2.104/24
utun0 172.168.23.87/24
```
