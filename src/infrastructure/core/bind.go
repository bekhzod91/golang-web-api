package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"io"
	"net/http"
	"strings"
)

func (c *context) BindJSON(target any) error {
	contentType := c.r.Header.Get("Content-Type")
	if contentType != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
		if mediaType != "application/json" {
			return errors.New("Content-Type header is not application/json")
		}
	}

	if c.r.ContentLength == 0 {
		return fmt.Errorf("body is empty")
	}

	c.r.Body = http.MaxBytesReader(c.w, c.r.Body, 1048576)
	defer c.r.Body.Close()

	data, err := io.ReadAll(c.r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}

func (c *context) ShouldBindJSON(target interface {
	Validate(formats strfmt.Registry) error
}) error {
	if err := c.BindJSON(target); err != nil {
		return err
	}

	return target.Validate(strfmt.Default)
}
