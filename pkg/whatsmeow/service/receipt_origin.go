package whatsmeow_service

import (
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
