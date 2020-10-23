package main

import (
    "fmt"
    ag_diagnostics "github.com/alec-lanter/grpc-diagnostics"
    "os"
)

func main() {
    dc, err := ag_diagnostics.ParseCommand(os.Args[1:])
    if err != nil {
        fmt.Println("ERROR: ", err)
        os.Exit(1)
    }

    result, err := dc.Execute()
    if err != nil {
        fmt.Println("ERROR: ", err)
        os.Exit(1)
    }

    fmt.Println("reply: ", result)
}
