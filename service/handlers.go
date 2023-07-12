package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) kickUser(c *gin.Context) {

	fmt.Println(c.Request.Body)
	var incomingReq ActiveFenceReq
	json.NewDecoder(c.Request.Body).Decode(&incomingReq)

	// Parse metadata
	var md ReqMetadata
	fmt.Println(incomingReq.Metadata)
	err := json.Unmarshal([]byte(incomingReq.Metadata), &md)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request: " + err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	channel := md.Cname

	// Call Agora API to kick
	url := "https://api.agora.io/dev/v1/kicking-rule"

	data := map[string]interface{}{
		"appid":           s.appID,
		"cname":           channel,
		"uid":             incomingReq.UID,
		"time_in_seconds": 300, // 5 minute ban
		"privileges":      []string{"join_channel"},
	}
	jsonData, _ := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Error: " + err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+s.restuflToken)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Request to Agora API Failed: " + err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	fmt.Println(string(body))

	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(`User %s kicked from channel %s`, incomingReq.UID, channel),
	})
}
