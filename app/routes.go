package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetEmails(store *EmailStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		pageInt, _ := strconv.Atoi(page)
		emails, err := store.ListEmails(pageInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, emails)
	}
}

func DeleteEmail(store *EmailStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, _ := strconv.Atoi(id)
		err := store.DeleteEmail(idInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Email eliminata correttamente"})
	}
}

func SendEmail(store *EmailStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()
		to := form.Value["recipient"][0]
		recipients := strings.Split(to, ",")
		email := Email{
			From:    form.Value["email"][0],
			To:      recipients,
			Subject: form.Value["subject"][0],
			Body:    form.Value["message"][0],
		}
		attachments := form.File["attachments"]
		err := SendEmailViaSMTP(email, attachments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Email inviata correttamente"})

	}
}
