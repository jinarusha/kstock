package httpclient

import (
	"bytes"
	"github.com/gojektech/heimdall/v6/httpclient"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

var (
	initOnce sync.Once
	client   *httpclient.Client
)

func GetClient() *httpclient.Client {
	initOnce.Do(func() {
		client = httpclient.NewClient()
	})

	return client
}

func EncodeJson(body interface{}) (*bytes.Buffer, error) {
	bytesData, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bytesData), nil
}
