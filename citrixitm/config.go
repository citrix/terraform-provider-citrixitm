package citrixitm

import (
	"log"
	"net/url"

	"github.com/cedexis/go-itm/itm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type config struct {
	ClientID     string
	ClientSecret string
	BaseURL      *url.URL
}

func newConfig(clientID string, clientSecret string, baseURL *url.URL) *config {
	result := config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL,
	}
	log.Printf("[DEBUG] Checking URL for trailing slash: %s", result.BaseURL.String())
	if 0 == len(result.BaseURL.Path) || "/" != result.BaseURL.Path[len(result.BaseURL.Path)-1:] {
		result.BaseURL.Path += "/"
	}
	return &result
}

// Returns an initialized itm.Client used to communicate with the Citrix ITM API
func (c *config) Client() (*itm.Client, error) {
	rel, _ := url.Parse("oauth/token")
	tokenURL := c.BaseURL.ResolveReference(rel)
	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     tokenURL.String(),
	}
	client, err := itm.NewClient(
		itm.HTTPClient(oauthConfig.Client(oauth2.NoContext)),
		itm.BaseURL(c.BaseURL))
	if err != nil {
		return nil, err
	}
	return client, nil
}
