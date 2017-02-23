package main

import (
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

func getUserPosts(session *fb.Session) <-chan FBResult {
	ch := make(chan FBResult)
	var paging *fb.PagingResult
	res, err := session.Get("/me/posts", fb.Params{"fields": "link,message,id"})
	if err != nil {
		// err can be an facebook API error.
		// if so, the Error struct contains error details.
		if e, ok := err.(*fb.Error); ok {
			log.Printf("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]",
				e.Message, e.Type, e.Code, e.ErrorSubcode)
		}
	} else {
		paging, err = res.Paging(session)
	}
	go func() {
		defer close(ch)
		noMore := false
		for !noMore {
			if err == nil {
				data := paging.Data()
				for i := range data {
					ch <- FBResult{post: FBUserPost{id: getStringFromMap(data[i], "id"), message: getStringFromMap(data[i], "message"), link: getStringFromMap(data[i], "link")}, err: nil}
				}
				noMore, err = paging.Next()
			} else {
				ch <- FBResult{post: *new(FBUserPost), err: err}
				break
			}
		}
	}()
	return ch
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

//FBResult User post or error
type FBResult struct {
	post FBUserPost
	err  error
}

//FBSettings FB API params
type FBSettings struct {
	appid string
	token string
}

//GetFBSettings read env and returns FBSettings
func GetFBSettings() *FBSettings {
	settings := new(FBSettings)
	settings.appid = os.Getenv(AppID)
	settings.token = os.Getenv(AppToken)
	return settings
}

var usertoken = flag.String("usertoken", "", "user token")

func main() {
	flag.Parse()
	settings := GetFBSettings()
	app := fb.New(settings.appid, settings.token)
	session := app.Session(*usertoken)
	for r := range getUserPosts(session) {
		if r.err != nil {
			log.Printf("got error %s", r.err)
		} else {
			log.Printf("result %s", r)
		}
	}
}
