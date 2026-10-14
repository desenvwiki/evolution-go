package whatsmeow_service

import (
    "fmt"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
)

// isRemoteRecipientReceipt checks the *source* of a receipt, not the message
// author or the provider's flattened "Delivered" state. Receipts from another
// linked device of our own account cannot prove recipient delivery.
func isRemoteRecipientReceipt(receipt *events.Receipt) bool {
    if receipt == nil || receipt.IsFromMe || receipt.Sender.IsEmpty() || receipt.Chat.IsEmpty() {
        return false
    }
    if receipt.IsGroup {
        return receipt.Chat.Server == types.GroupServer
    }
    return receipt.Chat.ToNonAD() == receipt.Sender.ToNonAD()
}

// The source belongs in the deduplication identity: a local-device echo
// must never consume the key used by the remote recipient's ACK.
func receiptDedupKey(instanceID, messageID, status string, sender types.JID) string {
    return fmt.Sprintf("%s_%s_%s_%s", instanceID, messageID, status, sender.ToNonAD().String())
}
