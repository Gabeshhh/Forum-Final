package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TABLE USERS
type Users struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

// TABLE POSTS
type Posts struct {
	gorm.Model
	TitlePost       string
	ContentCategory string
	UserID          string
	User            string
	Theme           string
	CommentsUser    []Comments `gorm:"foreignkey:PostID"`
	Links           string
	Date            time.Time
	PostsID         string `gorm:"column:PostsID"`

	// user_picture string
}

// TABLE COMMENTS
type Comments struct {
	gorm.Model
	Content string
	UserID  string
	User    Users
	PostID  string
	Post    Posts
}

func CreateDB(db *gorm.DB) {
	// Création des tables
	

	db.AutoMigrate(&Users{}, &Posts{}, &Comments{})
}

// Fonction pour ajouter un utilisateur

func AddUser(username, email, password string, db *gorm.DB) error {
	user := Users{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("Erreur lors de l'ajout de l'utilisateur:", err)
		return err
	}
	fmt.Println("Utilisateur ajouté avec succès :", user)
	return nil
}

// Fonction pour ajouter un post
func AddPost(title string, content string, theme string, db *gorm.DB) {
	// Générer un identifiant unique pour le post
	postID := uuid.NewString()

	// Ajout du post
	post := &Posts{
		TitlePost:       title,
		ContentCategory: content,
		Theme:           theme,
		PostsID:         postID,
		Date:            time.Now(),
	}

	db.Create(post)
}

// Fonction pour ajouter un commentaire
func AddComment(content string, userID string, postID string, db *gorm.DB) {
	// Ajout du commentaire
	db.Create(&Comments{Content: content, UserID: userID, PostID: postID})
}

func GetPostFromBdd(db *gorm.DB) ([]Posts, error) {
	var posts []Posts
	result := db.Find(&posts)

	for _, post := range posts {
		fmt.Println("Title:", post.TitlePost)
		fmt.Println("Content:", post.ContentCategory)
	}

	fmt.Println(result)

	return posts, nil
}

func GetAllPosts(db *gorm.DB) ([]Posts, error) {
	var posts []Posts
	if err := db.Find(&posts).Error; err != nil {
		return nil, err
	}

	for i := range posts {
		post := &posts[i]
		comments, err := GetCommentsByPostID(post.PostsID, db)
		if err != nil {
			return nil, err
		}
		post.CommentsUser = comments
	}
	return posts, nil
}

func GetCommentsByPostID(postID string, db *gorm.DB) ([]Comments, error) {
	var comments []Comments
	if err := db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
