package main

import (
	"fmt"
	"log"
	"net/http"

	database "Forum/Backend/Database"
	service "Forum/Backend/service"
	"Forum/Backend/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const IndexURL = "/index" // URL de la page d'accueil (utilise pour vérifier le token)

func main() {
	dsn := "forum.db"
	dbConnector, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}

	// Création et initialisation de la base de données

	database.CreateDB(dbConnector)

	// création route par default
	r := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.StaticFile("/styles.css", "./Frontend/styles.css") // permet de lier le fichier CSS
	r.StaticFile("/style.css", "./Frontend/style.css")
	r.LoadHTMLGlob("Frontend/*") // permet d'aller chercher les modèles HTML
	r.POST("/register", func(c *gin.Context) {
		service.RegisterUser(c, dbConnector)
		c.HTML(http.StatusOK, "log&Signup.html", gin.H{})
	})
	r.GET("/authentificationfield", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log&Signup.html", gin.H{})
	})

	// crée une route /get pour afficher la page
	r.GET("/dashboard-login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log&Signup.html", gin.H{})
	})

	// permet de vérifier les informations du méthod POST
	r.POST("/login", func(ctx *gin.Context) {
		loginService := service.StaticLoginService(dbConnector)
		var jwtService service.JWTService = service.JWTAuthService()
		var loginController controller.LoginController = controller.LoginHandler(loginService, jwtService)
		token := loginController.Login(ctx)
		if token != "" {
			fmt.Println("Token found")
			ctx.SetCookie("token", token, 3600, "/", "", false, false)
			ctx.Redirect(http.StatusFound, "/index")
		} else {
			fmt.Println("Token not found")
			ctx.Redirect(http.StatusFound, "/dashboard-login")
		}
	})

	// Middleware pour vérifier le token à chaque requête
	r.Use(func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			// Rediriger vers la page de connexion si le token est invalide ou non présent
			c.Redirect(http.StatusFound, "/authentificationfield")
			c.Abort()
			return
		}
		c.Next()
	})
	r.GET("/index", func(c *gin.Context) {
		// Récupérer tous les posts avec leurs commentaires associés
		posts, err := database.GetAllPosts(dbConnector)
		if err != nil {
			log.Println("Erreur lors de la récupération des posts:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Afficher les posts et leurs commentaires dans la table
		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts": posts,
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log&Signup.html", gin.H{})
	})

	r.GET("/create-post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create-post.html", gin.H{})
	})

	r.POST("/sendPost", func(c *gin.Context) {
		token, _ := c.Cookie("token")
		fmt.Println(token)
		service.CreatePost(c, dbConnector)
		c.Redirect(http.StatusFound, "/index")
	})

	r.POST("/comments", func(c *gin.Context) {
		postID := c.PostForm("PostsID")
		username := "Anonyme"
		content := c.PostForm("content")

		comment := database.Comments{
			Content: content,
			UserID:  username,
			PostID:  postID,
		}

		fmt.Println("Commentaire :", content)
		fmt.Println("PostID :", postID)

		// Appeler la fonction pour ajouter le commentaire dans la base de données
		database.AddComment(comment.Content, comment.UserID, comment.PostID, dbConnector)

		c.Redirect(http.StatusFound, "/index")
		// Rediriger vers la page d'accueil après l'ajout du commentaire
	})

	r.Run(":8089")
	fmt.Println("http://localhost:8089/dashboard-login")
}