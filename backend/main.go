package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"net/http"
	"strings"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Meeting struct {
	UUID      string    `json:"uuid"`
	ID        int       `json:"id"`
	HostID    string    `json:"host_id"`
	Topic     string    `json:"topic"`
	Type      int       `json:"type"`
	Duration  int       `json:"duration"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	StartURL  string    `json:"start_url"`
	JoinURL   string    `json:"join_url"`
	Settings  struct {
		HostVideo           bool   `json:"host_video"`
		ParticipantVideo    bool   `json:"participant_video"`
		CnMeeting           bool   `json:"cn_meeting"`
		InMeeting           bool   `json:"in_meeting"`
		JoinBeforeHost      bool   `json:"join_before_host"`
		MuteUponEntry       bool   `json:"mute_upon_entry"`
		Watermark           bool   `json:"watermark"`
		UsePmi              bool   `json:"use_pmi"`
		ApprovalType        int    `json:"approval_type"`
		Audio               string `json:"audio"`
		AutoRecording       string `json:"auto_recording"`
		EnforceLogin        bool   `json:"enforce_login"`
		EnforceLoginDomains string `json:"enforce_login_domains"`
		AlternativeHosts    string `json:"alternative_hosts"`
	} `json:"settings"`
}

func getJWT() (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": os.Getenv("ZOOM_API_KEY"),
		"exp": fmt.Sprintf("%d", time.Now().Add(60*time.Minute).Unix()),
	})
	
	tokenString, err = token.SignedString([]byte(os.Getenv("ZOOM_API_SECRET")))
	if err != nil {
		fmt.Println("Error Can not create JWT!")
		tokenString = ""
		return "", err
	}
	
	fmt.Println(tokenString)
	return tokenString, err
}

func createZoomMTG(jwt string) (url string) {
	payload := strings.NewReader(`{"type":1}`) // 1は通常の簡易ミーティング

	createMeetingUrl := "https://api.zoom.us/v2/users/" + os.Getenv("USER_ID") + "/meetings"

	req, err := http.NewRequest("POST", createMeetingUrl, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ jwt) // ヘッダに設定

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 201 {
		fmt.Println("status code is not 201!", res.Status)
	}

	defer res.Body.Close()
	var meeting Meeting
	err = json.NewDecoder(res.Body).Decode(&meeting)
	if err != nil {
		fmt.Println("json parse error!", err)
	}
	fmt.Println(meeting.JoinURL)
	url = meeting.JoinURL
	return url
}

func main() {
	err := godotenv.Load()
	if err != nil {
	  fmt.Println("Error loading .env file")
	}

	r := gin.Default()

	r.POST("/reserve/zoom/mtg", func(c *gin.Context) {
		jwt, err := getJWT()
		if err != nil {
			fmt.Println("Error Can not reserved Zoom MTG!")
			jwt = "can not reserved!"
		}
		url := createZoomMTG(jwt)
		c.JSON(200, gin.H{
			"message": "reserved",
			"zoom_mtg_url": url,
		})
	})
	r.Run()
}
