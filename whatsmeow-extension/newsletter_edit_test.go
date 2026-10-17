package whatsmeow

import (
	"encoding/json"
	"testing"
)

func TestNewsletterEditPayloadOnlyIncludesChangedFields(t *testing.T) {
	name := "Novo nome"
	vars, err := buildNewsletterEditVariables("120363000000000000@newsletter", NewsletterEditChanges{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(vars)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["newsletter_id"] != "120363000000000000@newsletter" {
		t.Fatalf("wrong jid: %v", decoded["newsletter_id"])
	}
	updates, ok := decoded["updates"].(map[string]any)
	if !ok {
		t.Fatalf("missing updates: %v", decoded)
	}
	if updates["name"] != name {
		t.Fatalf("wrong name: %v", updates)
	}
	if _, ok := updates["description"]; ok {
		t.Fatalf("unexpected description change: %v", updates)
	}
}

func TestNewsletterEditRejectsEmptyChanges(t *testing.T) {
	if _, err := buildNewsletterEditVariables("120363000000000000@newsletter", NewsletterEditChanges{}); err == nil {
		t.Fatal("expected empty update rejection")
	}
}
