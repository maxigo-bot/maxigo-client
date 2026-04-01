package maxigo

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// CheckPhoneNumbers checks which of the given phone numbers are registered in Max.
// Returns a list of phone numbers that exist. Phone numbers should be in
// international format without the "+" prefix (e.g., "79001234567").
// Corresponds to GET /notify/exists.
func (c *Client) CheckPhoneNumbers(ctx context.Context, phoneNumbers []string) ([]string, error) {
	q := make(url.Values)
	q.Set("phone_numbers", strings.Join(phoneNumbers, ","))

	var result checkPhoneNumbersResult
	if err := c.do(ctx, "CheckPhoneNumbers", http.MethodGet, "/notify/exists", q, nil, &result); err != nil {
		return nil, err
	}
	return result.ExistingPhoneNumbers, nil
}
