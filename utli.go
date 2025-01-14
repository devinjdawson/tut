package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs"
)

type follower struct {
	uid string
	//name       string
	followedAt string
}
type followed struct {
	uid string
	//name       string
	followingAt string
}

// User profile info
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Displayname     string `json:"displayname"`
	ProfileImageURL string `json:"profileImageURL"`
	FollowedAt      string `json:"followedAt"`
	UnfollowedAt    string `json:"unfollowedAt"`
}

// Unfollower user profile info
type Unfollower struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Displayname     string `json:"displayname"`
	ProfileImageURL string `json:"profileImageURL"`
	UnfollowedAt    string `json:"unfollowedAt"`
}

// Unfollowed user profile info
type Unfollowed struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Displayname     string `json:"displayname"`
	ProfileImageURL string `json:"profileImageURL"`
	UnfollowingAt   string `json:"unfollowedAt"`
}

// Notfollower user profile info
type Notfollower struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Displayname     string `json:"displayname"`
	ProfileImageURL string `json:"profileImageURL"`
	UnfollowedAt    string `json:"unfollowedAt"`
	UnfollowingAt   string `json:"unfollowedAt"`
}

type config struct {
	clientID       string
	oauth          string
	username       string
	userID         string
	serverPort     string
	updateInterval int
}

type apiResult struct {
	statusCode     int
	response       map[string]string
	limit          int
	limitRemaining int
	limtResetTime  int64
}

func getUserIDFromTwitch(username string, clientID string, oauth string) (apiResult, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username), nil)
	req.Header.Add("Client-ID", fmt.Sprintf("%s", clientID))
	if len(oauth) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauth))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	header := resp.Header
	limit, _ := strconv.Atoi(header["Ratelimit-Limit"][0])
	limitRemain, _ := strconv.Atoi(header["Ratelimit-Remaining"][0])
	limitReset, _ := strconv.ParseInt(header["Ratelimit-Reset"][0], 10, 64)

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		parsed, err := gabs.ParseJSON(body)
		if err != nil {
			log.Fatal(err)
		}

		exists := parsed.ExistsP("data.id")
		if exists {
			userdata, _ := parsed.Path("data").Children()
			firstuserdata, _ := userdata[0].ChildrenMap()
			return apiResult{resp.StatusCode,
				map[string]string{"id": firstuserdata["id"].Data().(string)},
				limit, limitRemain, limitReset}, nil
		}
	}

	return apiResult{resp.StatusCode, nil, limit, limitRemain, limitReset}, errors.New("getUserIDFromTwitch: cannot get UserID from Twitch API, check your ClientID or username")
}

func getUserNameFromTwitch(userID string, clientID string, oauth string) (apiResult, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?id=%s", userID), nil)
	req.Header.Add("Client-ID", fmt.Sprintf("%s", clientID))
	if len(oauth) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauth))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	header := resp.Header
	limit, _ := strconv.Atoi(header["Ratelimit-Limit"][0])
	limitRemain, _ := strconv.Atoi(header["Ratelimit-Remaining"][0])
	limitReset, _ := strconv.ParseInt(header["Ratelimit-Reset"][0], 10, 64)

	if resp.StatusCode == 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		parsed, err := gabs.ParseJSON(res)
		if err != nil {
			log.Fatal(err)
		}

		exists := parsed.ExistsP("data.login")
		if exists {
			userdata, _ := parsed.Path("data").Children()
			firstuserdata, _ := userdata[0].ChildrenMap()
			return apiResult{resp.StatusCode,
				map[string]string{"login": firstuserdata["login"].Data().(string), "displayname": firstuserdata["display_name"].Data().(string)},
				limit, limitRemain, limitReset}, nil
		}
	}

	return apiResult{resp.StatusCode, nil, limit, limitRemain, limitReset}, errors.New("getUserName: cannot get username from Twitch API, check your ClientID or userID")
}

func getUserFromTwitch(userID string, clientID string, oauth string) (apiResult, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?id=%s", userID), nil)
	req.Header.Add("Client-ID", fmt.Sprintf("%s", clientID))
	if len(oauth) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauth))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	header := resp.Header
	limit, _ := strconv.Atoi(header["Ratelimit-Limit"][0])
	limitRemain, _ := strconv.Atoi(header["Ratelimit-Remaining"][0])
	limitReset, _ := strconv.ParseInt(header["Ratelimit-Reset"][0], 10, 64)

	if resp.StatusCode == 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		parsed, err := gabs.ParseJSON(res)
		if err != nil {
			log.Fatal(err)
		}

		exists := parsed.ExistsP("data.login")
		if exists {
			userdata, _ := parsed.Path("data").Children()
			return apiResult{resp.StatusCode,
				map[string]string{"user": userdata[0].String()},
				limit, limitRemain, limitReset}, nil
		}
	}

	return apiResult{resp.StatusCode, nil, limit, limitRemain, limitReset}, errors.New("getUserName: cannot get username from Twitch API, check your ClientID or userID")
}

func getFollowersFromTwitch(userID string, pagination string, clientID string, oauth string) (apiResult, []follower, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users/follows?to_id=%s&first=100&after=%s", userID, pagination), nil)
	req.Header.Add("Client-ID", fmt.Sprintf("%s", clientID))
	if len(oauth) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauth))
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	header := resp.Header
	limit, _ := strconv.Atoi(header["Ratelimit-Limit"][0])
	limitRemain, _ := strconv.Atoi(header["Ratelimit-Remaining"][0])
	limitReset, _ := strconv.ParseInt(header["Ratelimit-Reset"][0], 10, 64)

	if resp.StatusCode == 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		parsed, err := gabs.ParseJSON(res)
		if err != nil {
			log.Fatal(err)
		}

		var output []follower
		followers, err := parsed.Path("data").Children()
		if err != nil {
			log.Fatal(err)
		}

		nextPagination := ""
		if parsed.Path("pagination.cursor").Data() != nil {
			nextPagination = parsed.Path("pagination.cursor").Data().(string)
		}

		if len(followers) > 0 {
			for _, child := range followers {
				childdata, _ := child.ChildrenMap()
				uid, _ := childdata["from_id"].Data().(string)
				followAt := childdata["followed_at"].Data().(string)
				output = append(output, follower{uid, followAt})
			}
		}
		return apiResult{resp.StatusCode,
			map[string]string{"next": nextPagination},
			limit, limitRemain, limitReset}, output, nil
	}
	return apiResult{resp.StatusCode, nil, limit, limitRemain, limitReset}, nil, errors.New("getFollowers: cannot get followers from Twitch API")
}

func getFollowingFromTwitch(userID string, pagination string, clientID string, oauth string) (apiResult, []followed, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users/follows?from_id=%s&first=100&after=%s", userID, pagination), nil)
	req.Header.Add("Client-ID", fmt.Sprintf("%s", clientID))
	if len(oauth) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauth))
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	header := resp.Header
	limit, _ := strconv.Atoi(header["Ratelimit-Limit"][0])
	limitRemain, _ := strconv.Atoi(header["Ratelimit-Remaining"][0])
	limitReset, _ := strconv.ParseInt(header["Ratelimit-Reset"][0], 10, 64)

	if resp.StatusCode == 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		parsed, err := gabs.ParseJSON(res)
		if err != nil {
			log.Fatal(err)
		}

		var output []followed
		following, err := parsed.Path("data").Children()
		if err != nil {
			log.Fatal(err)
		}

		nextPagination := ""
		if parsed.Path("pagination.cursor").Data() != nil {
			nextPagination = parsed.Path("pagination.cursor").Data().(string)
		}

		if len(following) > 0 {
			for _, child := range following {
				childdata, _ := child.ChildrenMap()
				uid, _ := childdata["to_id"].Data().(string)
				followAt := childdata["followed_at"].Data().(string)
				output = append(output, followed{uid, followAt})
			}
		}
		return apiResult{resp.StatusCode,
			map[string]string{"next": nextPagination},
			limit, limitRemain, limitReset}, output, nil
	}
	return apiResult{resp.StatusCode, nil, limit, limitRemain, limitReset}, nil, errors.New("getFollowing: cannot get following from Twitch API")
}
