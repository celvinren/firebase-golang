package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"restaurant_golang/models"
	"restaurant_golang/utils"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient  *storage.Client
	bucket         = "restaurant-flutter-c4308.appspot.com"
	credentialFile = "restaurant-flutter-firebase.json"
)

// InitFirebaseInstance init firebase app instance
func InitFirebaseInstance() (*firebase.App, error) {
	opt := option.WithCredentialsFile(credentialFile)
	app, errApp := firebase.NewApp(context.Background(), nil, opt)
	if errApp != nil {
		return nil, errApp
	}
	return app, errApp
}

// InitFireStore init firebase cloud store
func InitFireStore() (*firestore.Client, error) {
	app, errApp := InitFirebaseInstance()
	if errApp != nil {
		return nil, errApp
	}
	client, errClient := app.Firestore(context.Background())
	if errClient != nil {
		return nil, errClient
	}
	return client, errClient
}

// InitFirebaseAuth init firebase auth
func InitFirebaseAuth() (*auth.Client, error) {
	app, errApp := InitFirebaseInstance()
	if errApp != nil {
		return nil, errApp
	}
	client, errClient := app.Auth(context.Background())
	if errClient != nil {
		return nil, errClient
	}
	return client, errClient
}

// sendMessaging init firebase cloud store
func SendMessaging(c *gin.Context) {
	var msgNotificationObject models.Notification
	binderr := c.BindJSON(&msgNotificationObject)
	if binderr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "An error has occurred when bind notification:" + binderr.Error(),
			"data":   nil,
		})
		return
	}
	err := sendMessaging(msgNotificationObject)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "An error has occurred when bind notification: " + err.Error(),
			"data":   nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Successfully send notification",
		"data":   nil,
	})
}
func sendMessaging(notification models.Notification) error {
	// Create the message to be sent.
	msgNotification := &messaging.Message{
		Topic: notification.Topic,
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
		},
		Notification: &messaging.Notification{Body: notification.Body, Title: notification.Title},
	}

	app, errApp := InitFirebaseInstance()
	if errApp != nil {
		log.Println(errApp.Error())
		return errApp
	}
	msg, err := app.Messaging(context.Background())
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result, err := msg.Send(context.Background(), msgNotification)
	if err != nil {
		log.Println(err.Error())
		return err

	}
	log.Println(result)
	return nil
}

func SaveBase64ToCloudStore(base64Str string, filePath string, c *gin.Context) string {

	reader := utils.ConvertBase64StrToIOReader(base64Str)

	ctx := appengine.NewContext(c.Request)
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Println(err.Error())
	}

	sw := storageClient.Bucket(bucket).Object(filePath).NewWriter(ctx)
	if _, err := io.Copy(sw, reader); err != nil {
		log.Println("Could not copy file: ", err)
	}

	if err := sw.Close(); err != nil {
		log.Println("Could not close file: ", err)
	}

	u, _ := url.Parse("/" + bucket + "/" + sw.Attrs().Name)

	return "https://storage.googleapis.com" + u.EscapedPath()
}
