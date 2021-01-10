package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"restaurant_golang/models"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// RegisterUser godoc
// @Summary Register User
// @Description Register User
// @Success 200 {string} json
// @Router /apis/auth/registerUser [Post]
func RegisterUser(c *gin.Context) {
	var userDTO models.UserDTO
	bindErr := c.BindJSON(&userDTO)

	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "json cannot be explain: " + bindErr.Error(),
			"data":   nil,
		})
		return
	}

	authClient, err := InitFirebaseAuth()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "fail to create firebase client: " + bindErr.Error(),
			"data":   nil,
		})
		return
	}
	params := (&auth.UserToCreate{}).
		Email(userDTO.Email).
		Password(userDTO.Password)
	user, err := authClient.CreateUser(context.Background(), params)
	log.Printf("Successfully created user: %v\n", user)

	userDTO.Organization.UserID = user.UID
	userDTO.Password = ""
	isCreate, err := CreateOrganization(userDTO.Organization)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "fail to create create organization: " + bindErr.Error(),
			"data":   nil,
		})
	}
	if isCreate {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "Successfully create organization",
			"data":   models.UserDTO(userDTO),
		})
	}
}

// VerifyJWTToken verify email
func VerifyJWTToken(c *gin.Context) (bool, error) {
	result := false
	var err error
	authHeader := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, " ")
	reqUserID := splitToken[0]
	reqToken := splitToken[1]

	client, err := InitFirebaseAuth()
	if err != nil {
		log.Printf("Fail to verify your token")
	}
	token, err := client.VerifyIDToken(context.Background(), reqToken)
	if err != nil {
		log.Printf("Fail to verify your token")
	}
	if (time.Unix(token.Expires, 0).After(time.Now())) && (token.UID == reqUserID) {
		log.Printf("Verify JWT successfully")
		result = true
	} else {
		err = errors.New("Invalid token")
		log.Printf("Invalid token")
	}
	return result, err
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Success 200 {string} json
// @Router /apis/auth/verifyEmail [Get]
func VerifyEmail(c *gin.Context) {
	// mode := c.Request.URL.Query().Get("mode")
	oobCode := c.Request.URL.Query().Get("oobCode")

	body := make(map[string]string)
	body["oobCode"] = oobCode

	bytesData, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "convert request body error: %v" + err.Error(),
			"data":   nil,
		})
		return
	}
	fmt.Println(bytesData)

	reader := bytes.NewReader(bytesData)

	//get key from local file
	secretsContent, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		log.Printf(err.Error())
	}
	var secretsResult map[string]interface{}
	json.Unmarshal([]byte(secretsContent), &secretsResult)
	url := "https://identitytoolkit.googleapis.com/v1/accounts:update?key=" + secretsResult["firebaseKey"].(string)
	log.Printf(url)
	request, err := http.NewRequest("POST", url, reader)
	defer request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "request body close error: %v" + err.Error(),
			"data":   nil,
		})
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "call request error: %v" + err.Error(),
			"data":   nil,
		})
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "retrieve data from body error: %v" + err.Error(),
			"data":   nil,
		})
	}

	var response map[string]interface{}
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "retrieve data from body error: %v" + err.Error(),
			"data":   nil,
		})
	}

	statusCode := resp.StatusCode
	// str := (*string)(unsafe.Pointer(&respBytes))
	// fmt.Println("response body", *str, statusCode)
	if statusCode == 200 {
		log.Println(response["localId"].(string))
		//send notification
		notification := models.Notification{
			Topic: response["localId"].(string),
			Title: "Email Verifid",
			Body:  "Your email has been verifed!",
		}
		err = sendMessaging(notification)
		if err != nil {
			log.Println(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "Email verified",
			"data":   nil,
		})
	}

}
