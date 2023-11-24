package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/empfaze/golang_bot/utils"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	queryValues := url.Values{}

	queryValues.Add("offset", strconv.Itoa(offset))
	queryValues.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, queryValues)
	if err != nil {
		return nil, err
	}

	var result UpdatesResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "Can't do request"

	url := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, utils.WrapError(errMsg, err)
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, utils.WrapError(errMsg, err)
	}
	defer func() {
		response.Body.Close()
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, utils.WrapError(errMsg, err)
	}

	return body, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	const errMsg = "Can't send message"

	queryValues := url.Values{}
	queryValues.Add("chat_id", strconv.Itoa(chatID))
	queryValues.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, queryValues)
	if err != nil {
		return utils.WrapError(errMsg, err)
	}

	return nil
}
