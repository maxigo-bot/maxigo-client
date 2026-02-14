package maxigo

import (
	"context"
	"net/http"
	"net/url"
)

// Subscribe sets up a WebHook subscription for the bot.
// Corresponds to POST /subscriptions.
func (c *Client) Subscribe(ctx context.Context, webhookURL string, updateTypes []string, secret string) (*SimpleQueryResult, error) {
	body := SubscriptionRequestBody{
		URL:         webhookURL,
		UpdateTypes: updateTypes,
		Secret:      secret,
	}

	var result SimpleQueryResult
	if err := c.do(ctx, "Subscribe", http.MethodPost, "/subscriptions", nil, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Unsubscribe removes a WebHook subscription.
// After calling this method, long polling becomes available again.
// Corresponds to DELETE /subscriptions.
func (c *Client) Unsubscribe(ctx context.Context, webhookURL string) (*SimpleQueryResult, error) {
	q := make(url.Values)
	q.Set("url", webhookURL)

	var result SimpleQueryResult
	if err := c.do(ctx, "Unsubscribe", http.MethodDelete, "/subscriptions", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubscriptions returns all active WebHook subscriptions.
// Corresponds to GET /subscriptions.
func (c *Client) GetSubscriptions(ctx context.Context) ([]Subscription, error) {
	var result getSubscriptionsResult
	if err := c.do(ctx, "GetSubscriptions", http.MethodGet, "/subscriptions", nil, nil, &result); err != nil {
		return nil, err
	}
	return result.Subscriptions, nil
}
