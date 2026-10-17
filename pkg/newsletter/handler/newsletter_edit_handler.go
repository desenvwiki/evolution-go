package newsletter_handler

import (
    "errors"
    "net/http"

    instance_model "github.com/evolution-foundation/evolution-go/pkg/instance/model"
    newsletter_service "github.com/evolution-foundation/evolution-go/pkg/newsletter/service"
    "github.com/gin-gonic/gin"
)

// UpdateNewsletter updates only a channel owned/administered by the connected account.
func (n *newsletterHandler) UpdateNewsletter(ctx *gin.Context) {
    raw, exists := ctx.Get("instance")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error":"instance not found"})
        return
    }
    instance, ok := raw.(*instance_model.Instance)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid instance"})
        return
    }
    var input newsletter_service.UpdateNewsletterStruct
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    result, err := n.newsletterService.UpdateNewsletter(&input, instance)
    if err != nil {
        if errors.Is(err, newsletter_service.ErrNewsletterEditForbidden) {
            ctx.JSON(http.StatusForbidden, gin.H{"error":"only a channel administrator can edit this channel"})
            return
        }
        ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message":"success", "data":result})
}
