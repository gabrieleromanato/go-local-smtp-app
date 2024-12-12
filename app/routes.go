package app

import (
	"net/http"
	"strconv"

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
