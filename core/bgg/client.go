package bgg

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	Root = "https://www.boardgamegeek.com/xmlapi2"
)

type Client struct {
	Config ClientConfig
}

type ClientConfig struct {
	Root                string
	MaxRetries          int
	RetryDelayInSeconds int
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{Root: Root, MaxRetries: 10, RetryDelayInSeconds: 5}
}

func NewClient(config ClientConfig) *Client {
	return &Client{Config: config}
}

func (c *Client) GetThing(id string) (string, error) {
	url := c.Config.Root + "/thing?id=" + id

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		return string(body), nil
	} else {
		return string(body), errors.New("unexpected response code: " + strconv.Itoa(resp.StatusCode))
	}
}

func (c *Client) GetUser(name string) (*User, error) {
	url := c.Config.Root + "/users?name=" + name

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user User
	err = xml.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	if user.Id == "" {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (c *Client) GetCollection(filter CollectionFilter) (*Collection, error) {
	url := c.Config.Root + "/collection?username=" + filter.Username

	if filter.Subtype != "" {
		url += "&subtype=" + filter.Subtype
	}

	if filter.ExcludeSubtype != "" {
		url += "&excludesubtype=" + filter.ExcludeSubtype
	}

	if filter.Own {
		url += "&own=1"
	}

	if filter.Rated {
		url += "&rated=1"
	}

	if filter.Played {
		url += "&played=1"
	}

	if filter.Trade {
		url += "&trade=1"
	}

	if filter.Want {
		url += "&want=1"
	}

	if filter.Wishlist {
		url += "&wishlist=1"
	}

	if filter.Preordered {
		url += "&preordered=1"
	}

	if filter.PrevOwned {
		url += "&prevowned=1"
	}

	for i := 0; i < c.Config.MaxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusAccepted {
			time.Sleep(time.Duration(c.Config.RetryDelayInSeconds) * time.Second)
		} else if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			var collection Collection
			err = xml.Unmarshal(body, &collection)
			if err != nil {
				return nil, err
			}

			return &collection, nil
		} else {
			return nil, errors.New("unexpected response code: " + strconv.Itoa(resp.StatusCode))
		}
	}

	return nil, errors.New("max retries exceeded: " + strconv.Itoa(c.Config.MaxRetries))
}
