//go:generate protoc --go_out=. --go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative service.proto
package ag_diagnostics

import (
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
)

func getLocalIP() (string, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    for _, iface := range ifaces {
        if iface.Flags & (net.FlagUp | net.FlagLoopback) == 0 {
            // Interface down or loopback, keep looking
            continue
        }
        addrs, err := iface.Addrs()
        if err != nil {
            return "", err
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue // not something we're interested ini
            }
            ip = ip.To4()
            if ip == nil {
                continue // We don't want no V6 addresses!
            }
            return ip.String(), nil
        }
    }

    return "", errors.New("couldn't find an internet interface")
}

type DiagCommand struct {
    SourceIp string
    Endpoint string
    Command string
    Args []string
}

const usage = "usage : ag-diag <command> <host>:<port> [arguments]"

func ParseCommand(args []string) (DiagCommand, error) {
    if args[0] == "--help" {
        doHelp()
    }
    dc := DiagCommand{}
    if len(args) < 2 {
        return dc, fmt.Errorf("%s\n\nag-diag --help for help", usage)
    }
    dc.Command = args[0]
    endpoint := strings.Split(args[1], ":")
    if len(endpoint) != 2 {
        return dc, fmt.Errorf("bad endpoint address '%s': expected '<host>:<port>'", args[1])
    } else {
        dc.Endpoint = args[1]
    }

    if len(args) > 2 {
        dc.Args = args[2:]
    }

    ip, err := getLocalIP()
    if err != nil {
        return dc, err
    }

    dc.SourceIp = ip

    return dc, nil
}

func doHelp() {
    fmt.Println("\n", usage)
    fmt.Print("\n---- COMMANDS ----\n\n")
    for k, _ := range commands {
        fmt.Printf("  * %s\n", k)
    }
    fmt.Println()
    os.Exit(0)
}

type diagnosticFunction func(command DiagCommand) (string, error)

var commands = make(map[string]diagnosticFunction)

func (dc DiagCommand) Execute() (string, error) {
    f, ok := commands[dc.Command]
    if !ok {
        return "", errors.New(fmt.Sprintf("unrecognized command  '%s'", dc.Command))
    }

    return f(dc)
}
