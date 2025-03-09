package db

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"ppt/dao"
	"ppt/util"
)

type UserLoginLog struct {
	UserID    uint64 `json:"user_id"`
	UserName  string `json:"user_name"`
	LogType   string `json:"log_type"`
	IPAddress string `json:"ip_address"`
	Device    string `json:"device"`
	Timestamp int64  `json:"@timestamp"`
}

func GetUserLoginLog(userID uint64, specTime int64) ([]*UserLoginLog, error) {
	index := util.FormatESLogIndex(specTime)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"user_id":  userID,
				"log_type": "login",
			},
		},
		"sort": []map[string]interface{}{
			{"@timestamp": map[string]interface{}{
				"order": "desc",
			}},
		},
		"size": 10,
	}

	queryJson, err := json.Marshal(query)
	if err != nil {
		log.Printf("json marshal query err: %v", err)
		return nil, err
	}

	req := esapi.SearchRequest{
		Index:  []string{index},
		Body:   bytes.NewReader(queryJson),
		Pretty: true,
	}
	res, err := req.Do(context.Background(), dao.ESClient)
	if err != nil {
		log.Printf("es_api req_do error: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("json unmarshal res.Error err: %v", err)
			return nil, err
		}
		log.Fatalf("es search error: [%s] %s %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	var result map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("json unmarshal res.Body err: %v", err)
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	log.Printf("query found %d hits", len(hits))

	// 遍历命中
	loginLogs := make([]*UserLoginLog, len(hits))
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		data, err := json.Marshal(source)
		if err != nil {
			log.Printf("json marshal hits.source err: %v", err)
			continue
		}
		loginLog := &UserLoginLog{}
		err = json.Unmarshal(data, loginLog)
		if err != nil {
			log.Printf("json unmarshal hits.source.data err: %v", err)
			continue
		}
		loginLogs = append(loginLogs, loginLog)
	}
	return loginLogs, nil
}
