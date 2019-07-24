package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func backendServer(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", GetRoot).Methods("GET")
	router.HandleFunc("/followers", GetFollowers).Methods("GET")
	router.HandleFunc("/refollowers", GetRefollowers).Methods("GET")
	router.HandleFunc("/followersID", GetFollowersID).Methods("GET")
	router.HandleFunc("/unfollowers", GetUnfollowers).Methods("GET")
	router.HandleFunc("/following", GetFollowing).Methods("GET")
	router.HandleFunc("/refollowing", GetRefollowing).Methods("GET")
	router.HandleFunc("/followingID", GetFollowingID).Methods("GET")
	router.HandleFunc("/unfollowing", GetUnfollowing).Methods("GET")
	//	router.HandleFunc("/notfollowers", GetNonfollowers).Methods("GET")
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")
	fmt.Printf("[SYS] Server listening at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// GetRoot prints out root message
func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(`
	<p>Welcome to Twitch Unfollow Tracker</p>
	<p>Available Endpoints:</p>
	<ul>
	<li><a href="/followers">/followers</a></li>
	<li><a href="/refollowers">/refollowers</a></li>
	<li><a href="/followersID">/followersID</a></li>
	<li><a href="/unfollowers">/unfollowers</a></li>
	<li><a href="/following">/following</a></li>
	<li><a href="/refollowing">/refollowing</a></li>
	<li><a href="/followingID">/followingID</a></li>
	<li><a href="/unfollowing">/unfollowing</a></li>
	<li>a href="/notfollowers">/notfollowers</a></li>
	<li>/user/{id}</li>
	</ul>
	`))
}

// GetReFollowers find all refollowers detailed info
func GetRefollowers(w http.ResponseWriter, r *http.Request) {
	var outputUsers []User
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		uf := tx.Bucket([]byte("unfollowers"))
		f := tx.Bucket([]byte("followers"))
		u := tx.Bucket([]byte("users"))

		if f != nil && u != nil && uf != nil {
			uf.ForEach(func(k, v []byte) error {
				fdata := f.Get(k)
				if fdata == nil {
					return nil
				}
				udata := u.Get(k)
				if udata != nil && len(udata) > 0 {
					parsed, _ := gabs.ParseJSON(udata)
					jsondata, _ := parsed.ChildrenMap()
					out := User{
						string(k),
						jsondata["login"].Data().(string),
						jsondata["display_name"].Data().(string),
						jsondata["profile_image_url"].Data().(string),
						string(fdata),
						string(v)}
					outputUsers = append(outputUsers, out)
				} else {
					out := User{
						string(k),
						"",
						"",
						"",
						string(fdata),
						string(v)}
					outputUsers = append(outputUsers, out)
				}
				return nil
			})
		}
		return nil
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(outputUsers)
}

// GetReFollowing find all refollowing detailed info
func GetRefollowing(w http.ResponseWriter, r *http.Request) {
	var outputUsers []User
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		uo := tx.Bucket([]byte("unfollowing"))
		o := tx.Bucket([]byte("following"))
		u := tx.Bucket([]byte("users"))

		if o != nil && u != nil && uo != nil {
			uo.ForEach(func(k, v []byte) error {
				odata := o.Get(k)
				if odata == nil {
					return nil
				}
				udata := u.Get(k)
				if udata != nil && len(udata) > 0 {
					parsed, _ := gabs.ParseJSON(udata)
					jsondata, _ := parsed.ChildrenMap()
					out := User{
						string(k),
						jsondata["login"].Data().(string),
						jsondata["display_name"].Data().(string),
						jsondata["profile_image_url"].Data().(string),
						string(odata),
						string(v)}
					outputUsers = append(outputUsers, out)
				} else {
					out := User{
						string(k),
						"",
						"",
						"",
						string(odata),
						string(v)}
					outputUsers = append(outputUsers, out)
				}
				return nil
			})
		}
		return nil
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(outputUsers)
}

// GetFollowers find all followers detailed info
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	var outputUsers []User
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		f := tx.Bucket([]byte("followers"))
		uf := tx.Bucket([]byte("unfollowers"))
		u := tx.Bucket([]byte("users"))

		if f != nil && u != nil && uf != nil {
			f.ForEach(func(k, v []byte) error {
				udata := u.Get(k)
				ufdata := uf.Get(k)
				if udata != nil && len(udata) > 0 {
					parsed, _ := gabs.ParseJSON(udata)
					jsondata, _ := parsed.ChildrenMap()
					out := User{
						string(k),
						jsondata["login"].Data().(string),
						jsondata["display_name"].Data().(string),
						jsondata["profile_image_url"].Data().(string),
						string(v),
						string(ufdata)}
					outputUsers = append(outputUsers, out)
				} else {
					out := User{
						string(k),
						"",
						"",
						"",
						string(v),
						string(ufdata)}
					outputUsers = append(outputUsers, out)
				}
				return nil
			})
		}
		return nil
	})

	sort.Slice(outputUsers, func(i, j int) bool {
		return outputUsers[i].FollowedAt > outputUsers[j].FollowedAt
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(outputUsers)
}

// GetFollowing find all follows detailed info
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	var outputUsers []User
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		uo := tx.Bucket([]byte("unfollowing"))
		o := tx.Bucket([]byte("following"))
		u := tx.Bucket([]byte("users"))

		if o != nil && u != nil && uo != nil {
			o.ForEach(func(k, v []byte) error {
				udata := u.Get(k)
				uodata := uo.Get(k)
				if udata != nil && len(udata) > 0 {
					parsed, _ := gabs.ParseJSON(udata)
					jsondata, _ := parsed.ChildrenMap()
					out := User{
						string(k),
						jsondata["login"].Data().(string),
						jsondata["display_name"].Data().(string),
						jsondata["profile_image_url"].Data().(string),
						string(v),
						string(uodata)}
					outputUsers = append(outputUsers, out)
				} else {
					out := User{
						string(k),
						"",
						"",
						"",
						string(v),
						string(uodata)}
					outputUsers = append(outputUsers, out)
				}
				return nil
			})
		}
		return nil
	})

	sort.Slice(outputUsers, func(i, j int) bool {
		return outputUsers[i].FollowedAt > outputUsers[j].FollowedAt
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(outputUsers)
}

// GetFollowersID find all followers's ID
func GetFollowersID(w http.ResponseWriter, r *http.Request) {
	var followIDs []int
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		f := tx.Bucket([]byte("followers"))
		if f != nil {
			f.ForEach(func(k, v []byte) error {
				id, err := strconv.Atoi(string(k))
				if err != nil {
					log.Fatal(err)
				}
				followIDs = append(followIDs, id)
				return nil
			})
		}
		return nil
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(followIDs)
}

// GetFollowingID find all followers's ID
func GetFollowingID(w http.ResponseWriter, r *http.Request) {
	var followingIDs []int
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		o := tx.Bucket([]byte("following"))
		if o != nil {
			o.ForEach(func(k, v []byte) error {
				id, err := strconv.Atoi(string(k))
				if err != nil {
					log.Fatal(err)
				}
				followingIDs = append(followingIDs, id)
				return nil
			})
		}
		return nil
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(followingIDs)
}

// GetUnfollowers find all unfollowers
func GetUnfollowers(w http.ResponseWriter, r *http.Request) {
	var unfollowers []Unfollower
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		f := tx.Bucket([]byte("unfollowers"))
		f.ForEach(func(k, v []byte) error {
			u := tx.Bucket([]byte("users"))
			user := u.Get(k)
			parsed, err := gabs.ParseJSON(user)
			if err != nil {
				uf := Unfollower{
					string(k),
					"Unknown",
					"Unknown",
					"Unknown",
					string(v)}
				unfollowers = append(unfollowers, uf)
			} else {
				userdata, _ := parsed.ChildrenMap()
				uf := Unfollower{
					userdata["id"].Data().(string),
					userdata["login"].Data().(string),
					userdata["display_name"].Data().(string),
					userdata["profile_image_url"].Data().(string),
					string(v)}
				unfollowers = append(unfollowers, uf)
			}
			return nil
		})
		return nil
	})

	sort.Slice(unfollowers, func(i, j int) bool {
		return unfollowers[i].UnfollowedAt > unfollowers[j].UnfollowedAt
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(unfollowers)
}

// GetUnfollowing find all unfollowed
func GetUnfollowing(w http.ResponseWriter, r *http.Request) {
	var unfollowing []Unfollowed
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		o := tx.Bucket([]byte("unfollowing"))
		o.ForEach(func(k, v []byte) error {
			u := tx.Bucket([]byte("users"))
			user := u.Get(k)
			parsed, err := gabs.ParseJSON(user)
			if err != nil {
				uo := Unfollowed{
					string(k),
					"Unknown",
					"Unknown",
					"Unknown",
					string(v)}
				unfollowing = append(unfollowing, uo)
			} else {
				userdata, _ := parsed.ChildrenMap()
				uo := Unfollowed{
					userdata["id"].Data().(string),
					userdata["login"].Data().(string),
					userdata["display_name"].Data().(string),
					userdata["profile_image_url"].Data().(string),
					string(v)}
				unfollowing = append(unfollowing, uo)
			}
			return nil
		})
		return nil
	})

	sort.Slice(unfollowing, func(i, j int) bool {
		return unfollowing[i].UnfollowingAt > unfollowing[j].UnfollowingAt
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(unfollowing)
}

// GetUser get specific user
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	var user []byte
	db.View(func(tx *bolt.Tx) error {
		u := tx.Bucket([]byte("users"))
		user = u.Get([]byte(id))

		return nil
	})

	if user != nil {
		parsed, _ := gabs.ParseJSON(user)
		userdata, _ := parsed.ChildrenMap()
		uf := Unfollower{
			userdata["id"].Data().(string),
			userdata["login"].Data().(string),
			userdata["display_name"].Data().(string),
			userdata["profile_image_url"].Data().(string),
			""}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(uf)
	} else {
		w.WriteHeader(404)
	}
}
