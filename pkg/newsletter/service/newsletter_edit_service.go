package newsletter_service

import (
    "context"
    "errors"

    instance_model "github.com/evolution-foundation/evolution-go/pkg/instance/model"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
)

var ErrNewsletterEditForbidden = errors.New("only a channel owner or administrator can edit this channel")

type UpdateNewsletterStruct struct {
    JID           types.JID `json:"jid"`
    Name          *string   `json:"name,omitempty"`
    Description   *string   `json:"description,omitempty"`
    PictureBase64 string    `json:"picture_base64,omitempty"`
}

func canManageNewsletter(meta *types.NewsletterMetadata) bool {
    if meta == nil || meta.ViewerMeta == nil {
        return false
    }
    return meta.ViewerMeta.Role == types.NewsletterRoleOwner || meta.ViewerMeta.Role == types.NewsletterRoleAdmin
}

func (n *newsletterService) UpdateNewsletter(input *UpdateNewsletterStruct, instance *instance_model.Instance) (*types.NewsletterMetadata, error) {
    if input == nil || input.JID.Server != types.NewsletterServer {
        return nil, errors.New("invalid channel JID")
    }
    client, err := n.ensureClientConnected(instance.Id)
    if err != nil {
        return nil, err
    }
    // Recheck rights on the server; never trust frontend permission flags.
    metadata, err := client.GetNewsletterInfo(context.Background(), input.JID)
    if err != nil { return nil, err }
    if !canManageNewsletter(metadata) { return nil, ErrNewsletterEditForbidden }
    picture, err := decodeNewsletterPicture(input.PictureBase64)
    if err != nil { return nil, err }
    return client.UpdateNewsletter(context.Background(), input.JID, whatsmeow.NewsletterEditChanges{
        Name: input.Name, Description: input.Description, Picture: picture,
    })
}
