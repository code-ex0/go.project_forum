package message

import (
	"../../database"
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	CreateBy int    `json:"create_by"`
	PostBy   int    `json:"post_by"`
	Content  string `json:"content"`
}

/**
@brief: GetMessages display all messages
*/
func GetMessages() []Message {
	db := database.DBConn
	var posts []Message
	db.Find(&posts)
	return posts
}

/**
@brief: GetMessage get message by ID
*/
func GetMessage(id string) Message {
	db := database.DBConn
	var post Message
	db.Find(&post, id)
	return post
}

/**
@brief: GetMessagesFromUser get message by user ID
*/
func GetMessagesFromUser(idUser int) []Message {
	db := database.DBConn
	var posts []Message
	db.Where("create_by = ?", idUser).Find(&posts)
	return posts
}

/**
@brief: GetMessagesByPostId will display message from the post id
*/
func GetMessagesByPostId(idPost string) []Message {
	db := database.DBConn
	var posts []Message
	db.Where("post_by = ?", idPost).Find(&posts)
	return posts
}

/**
@brief: GetMessagesByPostId will display message from the user id
*/
func GetMessagesByUserId(id int) []Message {
	db := database.DBConn
	var posts []Message
	db.Where("create_by = ?", id).Find(&posts)
	return posts
}

/**
@brief: NewMessage create a new message in database
*/
func NewMessage(message *Message) Message {
	db := database.DBConn
	newMessage := Message{Content: message.Content, CreateBy: message.CreateBy, PostBy: message.PostBy}
	db.Create(&newMessage)
	return newMessage
}

/**
@brief DeleteMessage delete message by id in database
*/
func DeleteMessage(id string) bool {
	db := database.DBConn
	var message Message
	db.First(&message, id)
	if message.Content == "" {
		return false
	}
	db.Delete(&message)
	return true
}
