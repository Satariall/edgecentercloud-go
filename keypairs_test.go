package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyPairs_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		publicKey = "ssh-key"
	)

	keyPairs := []KeyPair{{PublicKey: publicKey}}
	URL := fmt.Sprintf("/v1/keypairs/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(keyPairs)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.KeyPairs.List(ctx)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, keyPairs) {
		t.Errorf("KeyPairs.List\n returned %+v,\n expected %+v", resp, keyPairs)
	}
}

func TestKeyPairs_Get(t *testing.T) {
	setup()
	defer teardown()

	keypair := &KeyPair{SSHKeyID: testResourceID}
	URL := fmt.Sprintf("/v1/keypairs/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(keypair)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.KeyPairs.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, keypair) {
		t.Errorf("KeyPairs.Get\n returned %+v,\n expected %+v", resp, keypair)
	}
}

func TestKeyPairs_Create(t *testing.T) {
	setup()
	defer teardown()

	keyPairCreateRequest := &KeyPairCreateRequest{
		SSHKeyName: "ssh-key",
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/keypairs/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(KeyPairCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, keyPairCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.KeyPairs.Create(ctx, keyPairCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestKeyPairs_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/keypairs/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.KeyPairs.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestKeyPairs_Share(t *testing.T) {
	setup()
	defer teardown()

	keyPairShareRequest := &KeyPairShareRequest{
		SharedInProject: true,
	}
	keypair := &KeyPair{SSHKeyID: testResourceID}
	URL := fmt.Sprintf("/v1/keypairs/%d/%d/%s/%s", projectID, regionID, testResourceID, keypairsSharePath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(KeyPairShareRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, keyPairShareRequest, reqBody)
		resp, _ := json.Marshal(keypair)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.KeyPairs.Share(ctx, testResourceID, keyPairShareRequest)
	require.NoError(t, err)

	assert.Equal(t, keypair, resp)
}
