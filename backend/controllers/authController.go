package controllers

import (
	model "IITK_Mail/models"
	"IITK_Mail/store"
	"IITK_Mail/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginHelper(c *gin.Context, mongoStore *store.MongoStore) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var userData model.User
	var isUser bool
	filter := bson.M{"email": user.Email}
	isUser, userData = mongoStore.IsUserExist(filter)

	if !isUser {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	err := VerifyPassword(user.Password, userData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	tokenString, err := token.GenerateToken(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the JWT token as a cookie
	// c.SetCookie("jwt", tokenString, int(time.Hour*24), "/", "", false, true)
	c.SetCookie("jwt", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
func UserHelper(c *gin.Context, mongoStore *store.MongoStore) {
	// Retrieve the JWT token from the cookie
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found in cookie"})
		return
	}

	// Log the retrieved JWT token
	log.Println("JWT token retrieved from cookie:", cookie)

	// Now you can pass the entire Gin context to the TokenValid function
	claims, err := token.TokenValid(c)
	if err != nil {
		// log.Println("ErrorTokenValid : ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token (TokenValid func)"})
		return
	}

	// If the token is valid, you can access the claims
	// If the token is valid, you can access all the claims
	var claimsMap = make(map[string]interface{})
	for key, value := range claims {
		claimsMap[key] = value
	}
	c.JSON(http.StatusOK, gin.H{"message": "User is authenticated", "claims": claimsMap})

	// userEmail := claims["useremail"].(string)
	// c.JSON(http.StatusOK, gin.H{"message": "User is authenticated", "userEmail": userEmail})
}

func RegisterHelper(c *gin.Context, mongoStore *store.MongoStore) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	user.Username = strings.TrimSpace(user.Username)
	err = mongoStore.InsertUserData(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert Data"})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func GetCookieHandler(c *gin.Context) {
	// cookie, err := c.Cookie("jwt")
	_, err := c.Cookie("jwt")

	if err != nil {
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	// c.String(http.StatusOK, "Cookie value: %s", cookie)
}
