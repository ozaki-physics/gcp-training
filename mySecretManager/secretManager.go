package mysecretmanager

import (
	"context"
	"fmt"
	"log"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetGCPSecretValue(projectId, name string, version int) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	fullName := "projects/" + projectId + "/secrets/" + name + "/versions/" + strconv.Itoa(version)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fullName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// ログに書き出す
	log.Printf("Plaintext: %s\n", string(result.Payload.Data))
	return string(result.Payload.Data), nil
}
