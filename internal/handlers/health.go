package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/firebase"
)

// HealthCheck checks the health of the service
func HealthCheck(log *logrus.Logger, firestoreService *firebase.FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("HealthCheck called")

		// Check Firestore connection
		if _, err := firestoreService.Client.Collection("health_check").Documents(context.Background()).GetAll(); err != nil {
			log.WithError(err).Error("Health check failed: Firestore connection is down")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"firestore": "down",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"firestore": "up",
		})
	}
}
