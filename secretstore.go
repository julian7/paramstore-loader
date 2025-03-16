package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type SecretStore struct {
	KeyID     string            `json:"key_id,omitempty"`
	ParamRoot string            `json:"param_root"`
	Secrets   map[string]Secret `json:"secrets"`
}

func readSecretStore(fname string, secretStore *SecretStore) error {
	fd, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("opening input file: %w", err)
	}

	defer fd.Close()
	dec := json.NewDecoder(fd)
	if err := dec.Decode(secretStore); err != nil {
		return fmt.Errorf("parsing input file: %w", err)
	}

	return nil
}

func (store *SecretStore) UpdateSecrets(ctx context.Context) error {
	ssmClient, err := newSSMClient(ctx, store.KeyID, store.ParamRoot)
	if err != nil {
		return fmt.Errorf("creating SSM client while updating secrets: %w", err)
	}
	for name, val := range store.Secrets {
		resp, err := ssmClient.get(ctx, name)
		if err != nil || string(val) != *resp.Parameter.Value {
			putResp, err := ssmClient.put(ctx, name, string(val))
			if err != nil {
				log.Printf("warning: could not write %s: %s", name, err.Error())
			} else {
				log.Printf("writing secret %s to tier %v", name, putResp.Tier)
			}
		} else {
			log.Printf("skipping secret %s", name)
		}
	}
	return nil
}
