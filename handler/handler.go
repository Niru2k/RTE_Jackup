package handler

import (
	//user defined packages
	"online/logs"
	"online/middleware"
	"online/models"
	"online/repository"

	//Inbuild packages
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	//Third party packages
	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

// This is for Signup
func (db Database) Signup(c echo.Context) error {
	var (
		data models.User
		role models.Roles
	)
	log := logs.Log()
	log.Info.Println("Message : 'signup-API called'")

	//Get user details from request body
	if err := c.Bind(&data); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.SignupReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": 400,
				"error":  stmt,
			})
		}
	}

	//validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(data.Email) {
		log.Error.Println("Error : 'Invalid Email' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "Invalid Email",
		})
	}

	//validate the password
	if len(data.Password) < 8 {
		log.Error.Println("Error : 'password must be greater than 8 characters' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "password must be greater than 8 characters",
		})
	}

	//validate the role
	if data.Role != "admin" && data.Role != "user" {
		log.Error.Println("Error : 'Invalid role' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "Invalid role",
		})
	}

	//To check if the user details already exist or not
	data, err := repository.ReadUserByEmail(db.Db, data)
	if err == nil {
		log.Error.Println("Error : 'user already exist' Status : 409")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 409,
			"error":  "user already exist",
		})
	}

	//To change the password into hashedPassword
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Printf("Error : '%s'\n", err)
		return nil
	}
	data.Password = string(password)

	//Select a role_id for specified role
	role, _ = repository.ReadRoleIdByRole(db.Db, data)
	data.RoleId = role.RoleId

	//Adding a user details into our database
	if err = repository.CreateUser(db.Db, data); err != nil {
		log.Error.Println("Error : 'email already exist' Status : 409")
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"status": 409,
			"error":  "email already exist",
		})
	}

	log.Info.Println("Message : 'signup successful!!!' Status : 200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    200,
		"message":   "signup successful!!!",
		"user data": data,
	})
}

// This is for Login
func (db Database) Login(c echo.Context) error {
	var data models.User
	log := logs.Log()
	log.Info.Println("Message : 'login-API called'")
	//Get mail-id and password from request body
	if err := c.Bind(&data); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.LoginReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Status": 400,
				"error":  stmt,
			})
		}
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(data.Email) {
		log.Error.Println("Error : 'Invalid Email' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "Invalid Email",
		})
	}

	//To verify if the user email is exist or not
	user, err := repository.ReadUserByEmail(db.Db, data)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err == nil {
			// Fetch a JWT token
			auth, err := repository.ReadTokenByUserId(db.Db, user)
			if err == nil {
				log.Info.Println("Message : 'login successful!!!' Status : 200")
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status":  200,
					"message": "Login Successful!!!",
					"token":   auth.Token,
				})
			}

			//Create a token
			token, err := middleware.CreateToken(user, c)
			if err != nil {
				return err
			}
			auth.UserId, auth.Token = user.UserId, token
			if err = repository.AddToken(db.Db, auth); err != nil {
				log.Error.Printf("Error : '%s' Status : 409\n", err)
				return c.JSON(http.StatusConflict, map[string]interface{}{
					"status": 409,
					"error":  err.Error(),
				})
			}

			log.Info.Println("Message : 'login successful!!!' Status : 200")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status":  200,
				"message": "Login Successful!!!",
				"token":   token,
			})
		}
		log.Error.Println("Error : 'incorrect password' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "incorrect password",
		})

	}
	log.Error.Println("Error : 'user not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "user not found",
	})
}

// Handler for post a product
func (db Database) PostProduct(c echo.Context) error {
	var Product models.ProductInfo
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'AddProduct-API called'")
	if err := c.Bind(&Product); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.ProductInfoReq{})
	for _, field := range fields {
		if reflect.ValueOf(&Product).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Status": 400,
				"error":  stmt,
			})
		}
	}
	if err := repository.CreateProduct(db.Db, Product); err != nil {
		log.Error.Printf("Error : '%s' Status : 400\n", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  err.Error(),
		})
	}
	log.Info.Println("Message : 'Product added successfully' Status : 200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Product added successfully",
	})
}

// Handler for update a product by product-id
func (db Database) UpdateProductById(c echo.Context) error {
	var check int
	log := logs.Log()

	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'UpdateProduct-API called'")
	Product, err := repository.ReadProductByProductId(db.Db, c.Param("product_id"))
	if err == nil {
		if err := c.Bind(&Product); err != nil {
			log.Error.Println("Error : 'internal server error' Status : 500")
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": 500,
				"error":  "internal server error",
			})
		}

		fields := structs.Names(models.ProductInfoReq{})
		for _, field := range fields {
			if reflect.ValueOf(&Product).Elem().FieldByName(field).Interface() == "" {
				check++
			}
		}
		if check == 4 {
			log.Error.Println("Error : 'no data found to do update' Status : 404")
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status": 404,
				"error":  "no data found to do update",
			})
		}
		if err := repository.UpdateProductByProductId(db.Db, c.Param("product_id"), Product); err == nil {
			log.Info.Println("Message : 'Product updated successfully' Status : 200")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status":  200,
				"message": "Product updated Successfully!!!",
			})
		}
	}
	log.Error.Println("Error : 'Product not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "Product not found",
	})
}
