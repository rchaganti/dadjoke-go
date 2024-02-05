package dj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	BaseUrl   *url.URL
	UserAgent string

	HttpClient *http.Client
}

func NewClient() *Client {
	baseUrl, _ := url.Parse("https://icanhazdadjoke.com")
	return &Client{
		BaseUrl:    baseUrl,
		UserAgent:  "dadjoke-go",
		HttpClient: http.DefaultClient,
	}
}

type Search struct {
	CurrentPage  int    `json:"current_page"`
	Limit        int    `json:"limit"`
	NextPage     int    `json:"next_page"`
	PreviousPage int    `json:"previous_page"`
	Results      []Joke `json:"results"`
	SearchTerm   string `json:"search_term"`
	Status       int    `json:"status"`
	TotalJokes   int    `json:"total_jokes"`
	TotalPages   int    `json:"total_pages"`
}

type Joke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

type SlackJoke struct {
	Attachments  []SlackAttachment `json:"attachments"`
	ResponseType string            `json:"response_type"`
	Username     string            `json:"username"`
}

type SlackAttachment struct {
	Fallback string `json:"fallback"`
	Footer   string `json:"footer"`
	Text     string `json:"text"`
}

func (c *Client) GetJoke() (Joke, error) {
	req, err := c.newRequest("GET", "", nil)
	if err != nil {
		return Joke{}, err
	}

	var joke Joke
	_, err = c.do(req, &joke, false)

	return joke, err
}

func (c *Client) GetJokeAsSlackMessage() (SlackJoke, error) {
	req, err := c.newRequest("GET", "/slack", nil)
	if err != nil {
		return SlackJoke{}, err
	}

	var slackJoke SlackJoke
	_, err = c.do(req, &slackJoke, false)

	return slackJoke, err
}

func (c *Client) GetJokeAsImage(id, imagePath string) error {
	path := fmt.Sprintf("/j/%s.png", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil, true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.Create(imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	return err
}

func (c *Client) SearchDadJokes(term string, page int, limit int) (Search, error) {
	search := fmt.Sprintf("/search?term=%s&page=%d&limit=%d", term, page, limit)
	req, err := c.newRequest("GET", search, nil)
	if err != nil {
		return Search{}, err
	}

	var searchResults Search
	_, err = c.do(req, &searchResults, false)

	return searchResults, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	t := c.BaseUrl.ResolveReference(rel)
	u, err := url.QueryUnescape(t.String())
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(r *http.Request, v interface{}, noDecode bool) (*http.Response, error) {
	resp, err := c.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if !noDecode {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
