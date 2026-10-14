package whatsmeow_service

import (
    "testing"

    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
)

func TestRemoteRecipientReceipt(t *testing.T) {
    remote := types.NewJID("5511999999999", types.DefaultUserServer)
    own := types.NewJID("5511888888888", types.DefaultUserServer)
    group := types.NewJID("120363123456789", types.GroupServer)

    cases := []struct {
        name string
        source types.MessageSource
        want bool
    }{
        {"remote recipient delivery", types.MessageSource{Chat:remote,Sender:remote},true},
        {"own linked-device echo", types.MessageSource{Chat:remote,Sender:own,IsFromMe:true},false},
        {"own echo with chat recipient", types.MessageSource{Chat:remote,Sender:own,IsFromMe:true},false},
        {"mismatched recipient", types.MessageSource{Chat:remote,Sender:own},false},
        {"remote group member", types.MessageSource{Chat:group,Sender:remote,IsGroup:true},true},
        {"own group device", types.MessageSource{Chat:group,Sender:own,IsGroup:true,IsFromMe:true},false},
        {"missing source",types.MessageSource{},false},
    }
    for _,tt := range cases {
        t.Run(tt.name,func(t *testing.T) {
            got := isRemoteRecipientReceipt(&events.Receipt{MessageSource:tt.source})
            if got != tt.want {t.Fatalf("got %v want %v",got,tt.want)}
        })
    }
}
