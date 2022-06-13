package controller

import (
	"../model/message"
	"../model/post"
	"../model/user"
	"../module"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

/**
 *	@brief: GetMessages get the message in the post
 */
func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages := message.GetMessages()
	messagesJson, _ := json.Marshal(messages)
	w.Write(messagesJson)
}

/**
 *	@brief: GetMessage get a message by his id
 */
func GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	gMessage := message.GetMessage(id)
	messageJson, _ := json.Marshal(gMessage)
	w.Write(messageJson)
}

/**
 * @brief: GetMessageByUserId work in progress
 */
func GetMessageByUserId(w http.ResponseWriter, r *http.Request) {
	User, _ := user.GetUserByEmail(r.Header.Get("Email"))
	PPost := message.GetMessagesByUserId(int(User.ID))
	postsJson, _ := json.Marshal(PPost)
	w.Write(postsJson)
}

/**
 *	@brief: GetMessageByPostId get message by the post ID they are related
 */
func GetMessageByPostId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	PPost := message.GetMessagesByPostId(id)
	postsJson, _ := json.Marshal(PPost)
	w.Write(postsJson)
}

/**
 *	@brief: NewMessage create a new message
 */
func NewMessage(w http.ResponseWriter, r *http.Request) {
	gMessage := new(message.Message)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&gMessage)

	if gMessage.Content != "" {
		if gMessage.PostBy != 0 && post.GetPost(strconv.Itoa(gMessage.PostBy)).Title != "" {
			User, _ := user.GetUserByEmail(r.Header.Get("Email"))
			gMessage.CreateBy = int(User.ID)
			newPost := message.NewMessage(gMessage)
			postsJson, _ := json.Marshal(newPost)
			w.Write(postsJson)
			return
		} else {
			var err module.Error
			err = module.SetError(err, "post not found")
			json.NewEncoder(w).Encode(err)
			return
		}

	}
	var err module.Error
	err = module.SetError(err, "content is null")
	json.NewEncoder(w).Encode(err)
}

/**
 *	@brief: DeleteMessage delete message
 */
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !message.DeleteMessage(id) {
		var err module.Error
		err = module.SetError(err, "Message not found")
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Write([]byte("Post successfully deleted"))
}
