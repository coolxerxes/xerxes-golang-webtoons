package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	emailEndpoint = "http://api.kopeechka.store/mailbox-get-message?full=1&id=%s&token=d8a01f21a050d8c94567fc8ca141694c&type=JSON&api=2.0"
	maxCount      = 30
)

var _ Mailer = (*Tidal)(nil)

type tidalApiResponse struct {
	Status string `json:"status"`
	Body   string `json:"fullmessage,omitempty"`
}

type Kopeechka struct {
	ID   string `json:"id"`
	Mail string `json:"mail"`
}
type Tidal struct {
	domain string
	client *http.Client
}

// NewTidalMailer return a Mailer using tidal.lol temp mail api
// please consider supporting the creator financially: https://t.me/modules
func NewTidalMailer(domain string) Mailer {

	t := &Tidal{
		domain: domain,
	}

	t.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    1024,
			MaxConnsPerHost: 100,
			IdleConnTimeout: 100 * time.Second,
		},
		Timeout: 100 * time.Second,
	}

	return t
}

func (t *Tidal) GetContent(id string) (string, error) {
	fmt.Printf("GETCONTENT CALLED ON ID: %s", id)
	slave := func() (*tidalApiResponse, error) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(emailEndpoint, id), nil)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}

		res, err := t.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do: %w", err)
		}
		defer res.Body.Close()
		bodyText, err2 := io.ReadAll(res.Body)
		if err2 != nil {
			log.Fatal(err2)
		}
		var response *tidalApiResponse
		fmt.Println(response)
		err_ := json.Unmarshal(bodyText, &response)
		if err_ != nil {
			return nil, err2
		}
		return response, nil
	}

	counter := 0
	for {
		if counter >= maxCount {
			return "", ErrNotFound
		}

		time.Sleep(1 * time.Second)
		result, err := slave()
		if err != nil {
			return "", fmt.Errorf("tidal: get: slave: %w", err)
		}
		if result.Status == "ERROR" {
			counter++
			time.Sleep(500 * time.Millisecond)
			continue
		}

		return result.Body, nil
	}
}

func (t *Tidal) RandomAddress() (string, string) {
	// return fmt.Sprintf("%s@%s", helpers.RandString(8), t.domain)
	httpp := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    1024,
			MaxConnsPerHost: 100,
			IdleConnTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.kopeechka.store/mailbox-get-email?api=2.0&spa=1&site=webtoon.com&sender=webtoon&regex=&mail_type=%s&token=d8a01f21a050d8c94567fc8ca141694c", t.domain), nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err1 := httpp.Do(req)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer res.Body.Close()
	bodyText, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var results Kopeechka
	errj := json.Unmarshal(bodyText, &results)
	if errj != nil {
		log.Fatal(errj)
	}
	return results.Mail, results.ID
}
