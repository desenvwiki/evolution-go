package whatsmeow

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

// NewsletterEditChanges contains only the metadata the administrator wants to change.
// A nil field means "leave unchanged"; an empty description is an explicit removal.
type NewsletterEditChanges struct {
	Name        *string
	Description *string
	Picture     []byte
}

func buildNewsletterEditVariables(jid string, changes NewsletterEditChanges) (map[string]any, error) {
	if !strings.HasSuffix(jid, "@newsletter") {
		return nil, errors.New("invalid newsletter JID")
	}
	updates := make(map[string]any)
	if changes.Name != nil {
		name := strings.TrimSpace(*changes.Name)
		if name == "" {
			return nil, errors.New("newsletter name cannot be empty")
		}
		updates["name"] = name
	}
	if changes.Description != nil {
		updates["description"] = *changes.Description
	}
	if len(changes.Picture) > 0 {
		updates["picture"] = changes.Picture
	}
	if len(updates) == 0 {
		return nil, errors.New("no newsletter changes provided")
	}
	return map[string]any{"newsletter_id": jid, "updates": updates}, nil
}

func (cli *Client) UpdateNewsletter(ctx context.Context, jid types.JID, changes NewsletterEditChanges) (*types.NewsletterMetadata, error) {
	variables, err := buildNewsletterEditVariables(jid.String(), changes)
	if err != nil {
		return nil, err
	}
	data, err := cli.sendMexIQ(ctx, mutationUpdateNewsletter, variables)
	if err != nil {
		return nil, err
	}
	var response struct {
		Newsletter *types.NewsletterMetadata `json:"xwa2_newsletter_update"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	if response.Newsletter == nil {
		return nil, errors.New("WhatsApp did not confirm newsletter update")
	}
	return response.Newsletter, nil
}
