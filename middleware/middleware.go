package middleware

import (
	//User-defined packages

	"online/helper"
	"online/logs"
	"online/models"
	"online/repository"

	//Inbuild packages
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	//Third-party packages
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Create a JWT token with the needed claims
func CreateToken(user models.User, c echo.Context) (string, error) {
	log := logs.Log()
	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}
	exp := time.Now().Add(time.Hour * 24).Unix()
	userId := strconv.Itoa(int(user.UserId))
	roleId := strconv.Itoa(int(user.RoleId))
	claims := jwt.MapClaims{
		"ExpiresAt": exp,
		"User-id":   userId,
		"IssuedAt":  time.Now().Unix(),
		"Role-id":   roleId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type Database struct {
	Db *gorm.DB
}

// Token and claims validation
func (db Database) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	log := logs.Log()
	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		//To check the token is empty or not
		if tokenString == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Status": 400,
				"Error":  "token is empty",
			})
		}

		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"Status": 400,
					"Error":  "Invalid token signature",
				})
			} else if claims["ExpiresAt"].(int64) < time.Now().Unix() {
				repository.DeleteToken(db.Db, claims["User-id"].(string))
				log.Error.Println("Error : 'session expired...login again!!!' Status : 440")
				return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
					"status": 440,
					"Error":  "session expired...login again!!!",
				})
			} else if !token.Valid {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"Status": 400,
					"Error":  "Invalid token",
				})
			}
		}

		// Check the user's role
		if claims["Role-id"] == "1" {
			c.Set("role", "admin")
		} else if claims["Role-id"] == "2" {
			c.Set("role", "user")
		}

		return next(c)
	}
}

// Get a claims from the token
func GetTokenClaims(c echo.Context) jwt.StandardClaims {
	log := logs.Log()
	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}
	tokenString := c.Request().Header.Get("Authorization")
	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}
	claims := jwt.StandardClaims{}
	jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	return claims
}

// Admin authorization
func AdminAuth(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return errors.New("unauthorized entry")
	}
	return nil
}

// User authorization
func UserAuth(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "user" {
		return errors.New("unauthorized entry")
	}
	return nil
}
