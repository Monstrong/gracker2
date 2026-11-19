package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Checker struct {
    client  *http.Client
    targets map[string]string
	timeout time.Duration
}

func NewChecker(targets map[string]string, clientTimeout time.Duration) *Checker {
    return &Checker{
        client: &http.Client{
            Timeout: clientTimeout,
        },
        targets: targets,
		timeout: clientTimeout,
    }
}

func (c *Checker) Check(ctx context.Context) error {
	for name, addr := range c.targets {
		req, err := http.NewRequestWithContext(
			ctx, 
			http.MethodGet, 
			"http://"+addr+"/ready", 
			nil,
		)
		if err != nil {
			return err
		}
		resp, err := c.client.Do(req)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s: returned status %d", name, resp.StatusCode)
		}
	}
	return nil
}

func (c *Checker) GetTimeout() time.Duration {
	return c.timeout
}