package dao

import (
	"crypto/tls"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"ppt/config"
	"sync"
)

var (
	ESClient *elasticsearch.Client
	esOnce   sync.Once
)

func InitESClient(cfg *config.ESConfig) error {
	esOnce.Do(func() {
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
			panic(err)
		}
		ping, err := client.Ping()
		if err != nil {
			panic(err)
		}
		if ping.IsError() {
			panic(ping)
		}
		ESClient = client
	})

	return nil
}
