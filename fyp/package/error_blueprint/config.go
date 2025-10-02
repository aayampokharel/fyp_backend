package errorz

import (
	"errors"
	"fmt"		
)

type Config struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *Config) Error() string {	
	return e.Message
}

func (c *Config) Wrap(m string) error {
	if c.Message == "" {
		return errors.New(m)
	}
	return fmt.Errorf("%w: %s", c, m)

}
