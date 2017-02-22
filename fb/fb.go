package main

import (
	"errors"
	"flag"
	"log"
	"os"

	fb "github.com/huandu/facebook"
)

const (
	// AppID FB app id
	AppID string = "APP_ID"
	//AppToken FB app token
	AppToken string = "APP_TOKEN"
)

func getUserPosts(usertoken string) ([]FBUserPost, error) {
	res, err := fb.Get("/me/posts", fb.Params{"access_token": usertoken, "fields": "link,message,id"})
	if err != nil {
		// err can be an facebook API error.
		// if so, the Error struct contains error details.
		if e, ok := err.(*fb.Error); ok {
			log.Printf("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]",
				e.Message, e.Type, e.Code, e.ErrorSubcode)
		}
	}
	data, ok := res["data"]
	if !ok {
		return nil, errors.New("data not found")
	}

	d, ok := data.([]interface{})
	if !ok {
		return nil, errors.New("Wrong data format")
	}

	r := make([]FBUserPost, len(d))
	for i := range d {
		v, ok := d[i].(map[string]interface{})
		if !ok {
			return nil, errors.New("Wrong data format")
		}
		r[i] = FBUserPost{link: getStringFromMap(v, "link"), message: getStringFromMap(v, "message"), id: getStringFromMap(v, "id")}
	}

	return r, err
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return ""
}

//FBUserPost User post
type FBUserPost struct {
	link    string
	message string
	id      string
}

//FBSettings FB API params
type FBSettings struct {
	appid string
	token string
}

//GetFBSettings read env and returns FBSettings
func GetFBSettings() (*FBSettings, error) {
	settings := new(FBSettings)
	settings.appid = os.Getenv(AppID)
	settings.token = os.Getenv(AppToken)
	return settings, nil
}

var usertoken = flag.String("usertoken", "", "user token")

func main() {
	flag.Parse()
	res, err := getUserPosts(*usertoken)
	if err != nil {
		log.Fatal(err)
	}

	for i := range res {
		log.Printf("result %s", res[i])
	}
}
