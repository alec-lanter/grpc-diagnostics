package ag_diagnostics

import (
    "context"
    "fmt"
    "github.com/alec-lanter/grpc-diagnostics/protobuf"
    "google.golang.org/grpc"
    "time"
)

func GetVersion(dc DiagCommand) (string, error) {
    conn, err := grpc.Dial(dc.Endpoint, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        return "", fmt.Errorf("could not connect: %w", err)
    }

    defer conn.Close()

    client := protobuf.NewDiagnosticServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    req := &protobuf.DiagnosticRequest{}
    req.SourceAddress = dc.SourceIp

    r, err := client.GetVersion(ctx, req)
    if err != nil {
        return "", fmt.Errorf("error invoking gRPC function: %w", err)
    }
    return r.Version, nil
}

