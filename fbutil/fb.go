package fbutil

import (
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

// GetPosts return channel of FBResults
func GetPosts(usertoken *string) <-chan FBResult {
	app := getApp()
	session := app.Session(*usertoken)
	return getUserPosts(session)
}

func getApp() *fb.App {
	settings := GetFBSettings()
	return fb.New(settings.appid, settings.token)
}

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
					ch <- FBResult{Post: FBUserPost{id: getStringFromMap(data[i], "id"), message: getStringFromMap(data[i], "message"), link: getStringFromMap(data[i], "link")}, Err: nil}
				}
				noMore, err = paging.Next()
			} else {
				ch <- FBResult{Post: *new(FBUserPost), Err: err}
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
	Post FBUserPost
	Err  error
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
