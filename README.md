# go-dns-playgroound

Play with Go DNS names resolution and take some notes about it.

## Go Name Resolution

Go resolves domain names using various methods that depends on the operating system.

Name resolution can be indirect (e.g. Dial function) or direct (e.g. LookupHost and LookupAddr functions).

### Unix

The resolver has two options for resolving names:

- Use a **pure Go** resolver that sends DNS requests directly to the server listed in the `/etc/resolv.conf` file
- Use a **cgo-based** resolver that calls C library routines such as getaddrinfo and getnameinfo

By default the pure Go resolver is used, because a blocked DNS request consumes only a goroutine, while a blocked C call consumes an operating system thread.

When a cgo is available, the cgo-based resolver is used instead under a variety of conditions:

- On systems that do not let programs to perform direct DNS requests (macOS)
- When the LOCALDOMAIN environment variable is present (even empty)
- When the RES_OPTIONS or HOSTALIASES environment variable is non-empty
- When the ASR_CONFIG environment variable is non-empty (Open-BSD only)
- When /etc/resolv.conf or /etc/nsswitch.conf specify the use of features that the Go resolver does not implement 
- When the name being looked up ends in .local or is a mDNS name

The resolver decision can be overridden by setting the netdns value of the GODEBUG environment variable

```bash
export GODEBUG=netdns=go    # force pure Go resolver
export GODEBUG=netdns=cgo   # force cgo resolver
```

The decision can also be forced while building the Go source tree by setting the `netgo` or `netcgo` build tag.

A numeric netdns setting, as in `GODEBUG=netdns=1`, causes the resolver to print debugging information about its decisions.

### Windows

On Windows, the resolver always uses C library functions, such as GetAddrInfo and DnsQuery.

## Run

### macOS

As described above, macOS doesn't let programs to perform direct DNS requests, so if you run the main function in the main.go file with the debug enabled

```bash
GODEBUG=netdns=1 go run dns.go
```

you will obtain the following output

```bash
❯ GODEBUG=netdns=1 go run dns.go
go package net: using cgo DNS resolver
2020/04/25 17:02:44 dial took 12.584126ms
2020/04/25 17:02:44 lookupHost took 554.512µs
2020/04/25 17:02:44 lookupAddr took 245.715µs
```

where the `net` package tells you that the `cgo DNS resolver` has been selected to resolve the domain name, as expected.
If you try to force the `cgo` resolver the program doesn't fail, but it takes a lot of time to perform the rresolution of the domain names (why?)

```bash
❯ GODEBUG=netdns=go go run dns.go
2020/04/25 17:12:30 dial took 54.249965ms
2020/04/25 17:12:30 lookupHost took 41.536915ms
2020/04/25 17:12:30 lookupAddr took 21.286505ms
```

### Linux

On Linux you can use both methods by setting the `GODEBUG=netdns` environment variable.

go

```bash
$ GODEBUG=netdns=go go run main.go
2020/04/25 16:58:24 dial took 2.553256ms
2020/04/25 16:58:24 lookupHost took 1.652838ms
2020/04/25 16:58:24 lookupAddr took 697.97µs
```

cgo

```bash
$ GODEBUG=netdns=go go run main.go
2020/04/25 16:58:30 dial took 3.19819ms
2020/04/25 16:58:30 lookupHost took 1.261135ms
2020/04/25 16:58:30 lookupAddr took 792.498µs
```

If you run your program without forcing the DNS resolver, the DNS resolver will be dynamically accorrding to the conditions mentioned aboove

```bash
go package net: dynamic selection of DNS resolver
2020/04/25 17:15:48 dial took 4.355174ms
2020/04/25 17:15:48 lookupHost took 2.970379ms
2020/04/25 17:15:48 lookupAddr took 662.399µs
```
