package newsletter_service

import (
    "testing"
    "go.mau.fi/whatsmeow/types"
)

func TestDecodeNewsletterPictureAcceptsPNGDataURL(t *testing.T) {
    raw := "data:image/png;base64,iVBORw0KGgo="
    value, err := decodeNewsletterPicture(raw)
    if err != nil || len(value) == 0 { t.Fatalf("expected decoded PNG, bytes=%d err=%v", len(value), err) }
}

func TestDecodeNewsletterPictureRejectsUnsupportedDataURL(t *testing.T) {
    if _, err := decodeNewsletterPicture("data:text/html;base64,PHNjcmlwdD4="); err == nil {
        t.Fatal("must reject non-image input")
    }
}

func TestCanManageNewsletterRole(t *testing.T) {
    tests := []struct{ role types.NewsletterRole; allowed bool }{
        {types.NewsletterRoleOwner, true},
        {types.NewsletterRoleAdmin, true},
        {types.NewsletterRoleSubscriber, false},
        {types.NewsletterRoleGuest, false},
    }
    for _, tc := range tests {
        got := canManageNewsletter(&types.NewsletterMetadata{
            ViewerMeta: &types.NewsletterViewerMetadata{Role: tc.role},
        })
        if got != tc.allowed { t.Errorf("role %s allowed=%v, want %v", tc.role, got, tc.allowed) }
    }
    if canManageNewsletter(&types.NewsletterMetadata{}) { t.Fatal("missing viewer role must fail closed") }
    if canManageNewsletter(nil) { t.Fatal("missing newsletter must fail closed") }
}
