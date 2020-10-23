package ag_diagnostics

import (
    "context"
    "errors"
    "fmt"
    "ag_diagnostics/protobuf"
    "google.golang.org/grpc"
    "time"
)

func GetVersion(dc DiagCommand) (string, error) {
    conCtx, cancel := context.WithTimeout(context.Background(), time.Second)

    conn, err := grpc.DialContext(conCtx, dc.Endpoint, grpc.WithInsecure())
    select {
        case _, ok := <-conCtx.Done():
            if !ok {
                // Channel is closed, our work is finished
                break
            }
            return "", errors.New("timeout connecting to host")
    }

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

