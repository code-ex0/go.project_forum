package post

import (
	"../../database"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

// create Post structure
type Post struct {
	gorm.Model
	CreateBy int    `json:"create_by"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Like     string `json:"-"`
	DisLike  string `json:"-"`
}

// create ResponsePost (for answers in a topic)
type ResponsePost struct {
	ID        uint      `json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Like      int       `json:"like"`
	DisLike   int       `json:"dislike"`
}

/**
@brief: GetPosts display a list of posts
*/
func GetPosts() []Post {
	db := database.DBConn
	var posts []Post
	db.Find(&posts)
	return posts
}

/**
@brief: GetPost call a precise post by his ID
*/
func GetPost(id string) Post {
	db := database.DBConn
	var post Post
	db.Find(&post, id)
	return post
}

/**
@brief: GetPostsByUserId will show a user post list
*/
func GetPostsByUserId(idUser string) []Post {
	db := database.DBConn
	var posts []Post
	db.Where("create_by = ?", idUser).Find(&posts)
	return posts
}

/**
@brief: NewPost create a post in the database
*/
func NewPost(post *Post) Post {
	db := database.DBConn
	newPost := Post{Title: post.Title, Content: post.Content, CreateBy: post.CreateBy}
	db.Create(&newPost)
	return newPost
}

/**
@brief: DeletePost delete a post
*/
func DeletePost(id string) bool {
	db := database.DBConn
	var post Post
	db.First(&post, id)
	if post.Title == "" {
		return false
	}
	db.Delete(&post)
	return true
}

/**
@brief: Like will add a like on a post (will add address email i a list (email a unique so 1 email = 1 user))
*/
func Like(Email string, id string) {
	db := database.DBConn // open database
	post := GetPost(id)   // get the post by his ID
	if post.Like != "" {  // check if there is a like
		found := false
		for _, v := range strings.Split(post.Like, ";") { // will separate mail by a ";" to make an array
			if Email == v {
				found = true
			}
		}
		if !found {
			post.Like = strings.Join([]string{post.Like, Email}, ";") // if the mail isn't in it already will add a like
		} else {
			post.Like = removeString(strings.Split(post.Like, ";"), Email) // if the mail is in it will remove the like
		}
	} else {
		post.Like = Email // if there is no like already add the first email into the list
	}
	if post.DisLike != "" { // auto remove dislike if the user like
		post.DisLike = removeString(strings.Split(post.DisLike, ";"), Email)
	}
	db.Save(&post)
}

func DisLike(Email string, id string) {
	db := database.DBConn
	post := GetPost(id)
	if post.DisLike != "" {
		found := false
		for _, v := range strings.Split(post.DisLike, ";") {
			if Email == v {
				found = true
			}
		}
		if !found {
			post.DisLike = strings.Join([]string{post.DisLike, Email}, ";")
		} else {
			post.DisLike = removeString(strings.Split(post.DisLike, ";"), Email)
		}
	} else {
		post.DisLike = Email
	}
	if post.Like != "" {
		post.Like = removeString(strings.Split(post.Like, ";"), Email)
	}
	db.Save(&post)
}

/**
@brief: removeString remove an element in a array of string
*/
func removeString(list []string, element string) string {
	temp := ""
	for _, v := range list {
		if element != v {
			if temp == "" {
				temp = v
			} else {
				strings.Join([]string{temp, v}, ";")
			}
		}
	}
	return temp
}
