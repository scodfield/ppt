package dao

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"ppt/config"
)

var (
	ESClient *elasticsearch.Client
)

func InitESClient(cfg *config.ESConfig) error {
	var err error
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{address},
		Username:  cfg.UserName,
		Password:  cfg.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		},
	})
	if err != nil {
		return err
	}
	ping, err := client.Ping()
	if err != nil {
		return err
	}
	if ping.IsError() {
		return errors.New("es client ping error")
	}
	ESClient = client
	return nil
}
