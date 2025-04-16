package service

import (
	"fmt"

	database "Forum/Backend/Database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context, db *gorm.DB) { // token string, db *gorm.DB) { pour l'username plus tard

	fmt.Println("____________________CREATE_______________________POST_______________________")

	title := c.PostForm("postUserTitle")
	content := c.PostForm("postUserContent")
	Theme := c.PostForm("spec")

	post := database.Posts{
		TitlePost:       title,
		ContentCategory: content,
		Theme:           Theme,
	}

	// Utiliser postId dans votre logique ou pour l'affecter à une autre variable si nécessaire

	// fmt.Println(username)
	fmt.Println(post)

	database.AddPost(post.TitlePost, post.ContentCategory, post.Theme, db)
}
