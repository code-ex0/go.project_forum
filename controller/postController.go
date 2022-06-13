package controller

import (
	"../model/post"
	"../model/user"
	"../module"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

/**
 * @brief: GetPosts display all posts
 */
func GetPosts(w http.ResponseWriter, r *http.Request) {
	posts := post.GetPosts()
	postsJson, _ := json.Marshal(posts)
	w.Write(postsJson)
}

/**
 *	@brief: GetPost select a post by his ID
 */
func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	PPost := post.GetPost(id)
	RPost := new(post.ResponsePost)

	RPost.Author = user.GetUserById(PPost.CreateBy).Pseudo
	RPost.CreatedAt = PPost.CreatedAt
	RPost.Title = PPost.Title
	RPost.Content = PPost.Content
	RPost.ID = PPost.ID
	if PPost.Like != "" {
		RPost.Like = len(strings.Split(PPost.Like, ";"))
	} else {
		RPost.Like = 0
	}
	if PPost.DisLike != "" {
		RPost.DisLike = len(strings.Split(PPost.DisLike, ";"))
	} else {
		RPost.DisLike = 0
	}
	postsJson, _ := json.Marshal(RPost)
	w.Write(postsJson)
}

/**
 * @brief: GetPostByUserId get every post for an user from his ID
 */
func GetPostByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	PPost := post.GetPostsByUserId(id)
	postsJson, _ := json.Marshal(PPost)
	w.Write(postsJson)
}

/**
 *	@brief: create a new post
 */
func NewPost(w http.ResponseWriter, r *http.Request) {
	PPost := new(post.Post)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&PPost)

	if PPost.Title != "" && PPost.Content != "" && len(PPost.Title) > 5 {
		User, _ := user.GetUserByEmail(r.Header.Get("Email"))
		PPost.CreateBy = int(User.ID)
		newPost := post.NewPost(PPost)
		postsJson, _ := json.Marshal(newPost)
		w.Write(postsJson)
		return
	}
	var err module.Error
	err = module.SetError(err, "content is null")
	json.NewEncoder(w).Encode(err)
	return
}

/**
 * @brief: DeletePost delete a post using his ID
 */
func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !post.DeletePost(id) {
		w.Write([]byte("Post not found"))
		return
	}

	w.Write([]byte("Post successfully deleted"))
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	post.Like(r.Header.Get("Email"), id)
}

func DisLikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	post.DisLike(r.Header.Get("Email"), id)
}
