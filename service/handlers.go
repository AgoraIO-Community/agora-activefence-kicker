package kickService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// KickHandler handles incoming moderation kick requests from ActiveFence.
// It parses the request, calls the Agora API to ban the user,
// and returns the response.
func (s *Service) KickHandler(w http.ResponseWriter, r *http.Request) {

	// Decode JSON request body into struct
	var incomingReq ActiveFenceReq
	if err := json.NewDecoder(r.Body).Decode(&incomingReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse channel name from request metadata
	var md ReqMetadata
	if err := json.Unmarshal([]byte(incomingReq.Metadata), &md); err != nil {
		http.Error(w, "Invalid metadata: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call service method to ban user
	// 300 seconds/5 minutes duration
	if err := KickUser(s.appID, md.Cname, incomingReq.UID, 300, s.restToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User %s kicked from channel %s", incomingReq.UID, md.Cname),
	})
}

// KickUser calls the Agora API to ban a user from a channel.
//
// It takes in the app ID, channel name, user ID, ban duration in seconds,
// and Agora REST API token.
//
// The ban removes the "join_channel" privilege from the user for the
// specified duration.
//
// Returns any error encountered, or nil if successful.
func KickUser(appId string, channel string, userId string, duration int, restToken string) error {
	// Call Agora API to kick

	url := "https://api.agora.io/dev/v1/kicking-rule"

	data := map[string]interface{}{
		"appid":           appId,
		"cname":           channel,
		"uid":             userId,
		"time_in_seconds": duration,
		"privileges":      []string{"join_channel"},
	}

	jsonData, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return fmt.Errorf("internal Error: %s", err.Error())
	}

	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+restToken)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("request to Agora API Failed: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(string(body))
	return nil
}
