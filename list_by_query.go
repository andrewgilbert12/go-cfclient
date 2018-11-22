package cfclient

import (
	"net/url"
	"github.com/pkg/errors"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type NextUrlResponse struct {
	NextUrl string `json:"next_url"`
}

func (c *Client) ListByQuery(errorMsg string, endpoint string, query url.Values) ([]json.RawMessage, error) {
	var results []json.RawMessage
	requestUrl := fmt.Sprintf("%s?%s",
		endpoint,
		query.Encode(),
	)

	_, listSinglePage := query["page"]

	for {
		var nextUrlResp NextUrlResponse

		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrapf(err, "Error requesting %s", errorMsg)
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading %s request", errorMsg)
		}
		results = append(results, resBody)

		json.Unmarshal(resBody, &nextUrlResp)
		requestUrl = nextUrlResp.NextUrl
		if listSinglePage || requestUrl == "" {
			break
		}
	}

	return results, nil
}
