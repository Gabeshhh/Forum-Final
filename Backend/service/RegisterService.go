package service

import (
	"net/http"

	database "Forum/Backend/Database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context, db *gorm.DB) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Vérifie si l'utilisateur existe déjà
	var existingUser database.Users
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusOK, "log&Signup.html", gin.H{
			"error": "Email déjà utilisé",
		})
		return
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "log&Signup.html", gin.H{
			"error": "Erreur de hash du mot de passe",
		})
		return
	}

	// Ajout en base
	if err := database.AddUser(username, email, string(hashedPassword), db); err != nil {
		c.HTML(http.StatusInternalServerError, "log&Signup.html", gin.H{
			"error": "Erreur lors de l'inscription",
		})
		return
	}

	c.HTML(http.StatusOK, "log&Signup.html", gin.H{
		"success": "Inscription réussie ! Tu peux te connecter 🚀",
	})
}
