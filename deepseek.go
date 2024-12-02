package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DeepSeek struct {
	Token string
	Proxy string
}

func (ds *DeepSeek) QueryString(userQuery, systemQuery string) (string, error) {
	return ds.Query([]Message{{Role: "assistant", Content: systemQuery}, {Role: "user", Content: userQuery}})
}

func (ds *DeepSeek) Query(messages []Message) (string, error) {

	request := Request{
		Model:    "deepseek-chat",
		Messages: messages,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ds.Token))

	client := &http.Client{}
	if ds.Proxy != "" {
		proxyURL, err2 := url.Parse(ds.Proxy)
		if err2 != nil {
			return "", err2
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", err
		}
		return response.Choices[0].Message.Content, nil
	}

	return "", errors.New(fmt.Sprintf("Error: %d", resp.StatusCode))

}
