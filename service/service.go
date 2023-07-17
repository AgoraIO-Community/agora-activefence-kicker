package kickService

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

// Service is the backend service
type Service struct {
	Server    *http.Server
	Sigint    chan os.Signal
	appID     string
	restToken string
}

// Stop service safely, closing additional connections if needed.
func (s *Service) Stop() {
	// Will continue once an interrupt has occurred
	signal.Notify(s.Sigint, os.Interrupt)
	<-s.Sigint

	// cancel would be useful if we had to close third party connection first
	// Like connections to a db or cache
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
}

// Start runs the service by listening to the specified port
func (s *Service) Start() {
	log.Println("Listening to port " + s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func NewService() *Service {

	godotenv.Load()

	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	customerKey, customerKeyExists := os.LookupEnv("CUSTOMER_KEY")
	customerSecret, customerSecretExists := os.LookupEnv("CUSTOMER_SECRET")
	serverPort, serverPortExists := os.LookupEnv("SERVER_PORT")
	if !appIDExists || !customerKeyExists || !customerSecretExists ||
		len(appIDEnv) == 0 || len(customerKey) == 0 || len(customerSecret) == 0 {
		log.Fatal("FATAL ERROR: ENV not properly configured, check environment variables APP_ID and CUSTOMER_KEY and CUSTOMER_SECRET")
	}
	if !serverPortExists || len(serverPort) == 0 {
		// Check $PORT, this is used by Railway.
		port, portExists := os.LookupEnv("PORT")
		if portExists && len(port) > 0 {
			serverPort = port
		} else {
			serverPort = "8080"
		}
	}

	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	s := &Service{
		Sigint: make(chan os.Signal, 1),
		Server: &http.Server{
			Addr: fmt.Sprintf(":%s", serverPort),
		},
		appID:     appIDEnv,
		restToken: base64Credentials,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/kick/", s.KickHandler)
	mux.HandleFunc("/kick", s.KickHandler)
	s.Server.Handler = mux
	return s
}

type ActiveFenceReq struct {
	UID      string `json:"userId"`
	Metadata string `json:"metadata"`
}
type ReqMetadata struct {
	Cname string `json:"cname"`
	UID   int    `json:"uid"`
	// other fields
}
