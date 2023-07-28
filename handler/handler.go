package handler

import (
	//user defined packages
	"online/logs"
	"online/middleware"
	"online/models"
	"online/repository"
	"strconv"

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
	Connection *gorm.DB
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
	data, err := repository.ReadUserByEmail(db.Connection, data)
	if err == nil {
		log.Error.Println("Error : 'user already exist' Status : 400")
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"status": 400,
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
	role, _ = repository.ReadRoleIdByRole(db.Connection, data)
	data.RoleId = role.RoleId

	//Adding a user details into our database
	if err = repository.CreateUser(db.Connection, data); err != nil {
		log.Error.Println("Error : 'email already exist' Status : 400")
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"status": 400,
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
	user, err := repository.ReadUserByEmail(db.Connection, data)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err == nil {
			// Fetch a JWT token
			auth, err := repository.ReadTokenByUserId(db.Connection, user)
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
			if err = repository.AddToken(db.Connection, auth); err != nil {
				log.Error.Printf("Error : '%s' Status : 400\n", err)
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"status": 400,
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
			"status": 403,
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
	if err := repository.CreateProduct(db.Connection, Product); err != nil {
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

// Handler for get all products
func (db Database) GetAllProducts(c echo.Context) error {
	log := logs.Log()
	log.Info.Println("Message : 'GetAllProducts-API called'")
	Products, err := repository.ReadAllProducts(db.Connection)
	if err == nil {
		log.Info.Println("Message : 'Product(s) retrieved successfully' Status : 200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":   200,
			"Products": Products,
		})
	}
	log.Error.Println("Error : 'Product not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "Product not found",
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
	Product, err := repository.ReadProductByProductId(db.Connection, c.Param("product_id"))
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
		if err := repository.UpdateProductByProductId(db.Connection, c.Param("product_id"), Product); err == nil {
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

// Handler for delete a product by product-id
func (db Database) DeleteProductById(c echo.Context) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Deleteproduct-API called'")
	if _, err := repository.ReadProductByProductId(db.Connection, c.Param("product_id")); err == nil {
		repository.DeleteProductByProductId(db.Connection, c.Param("product_id"))
		log.Info.Println("Message : 'Product deleted successfully' Status : 200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  200,
			"message": "Product deleted Successfully!!!",
		})
	}

	log.Error.Println("Error : 'Product not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "Product not found",
	})
}

// Handler for post a order
func (db Database) AddOrder(c echo.Context) error {
	var order models.OrderProductInfo
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'AddOrder-API called'")
	if err := c.Bind(&order); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.OrderProductReq{})
	for _, field := range fields {
		if reflect.ValueOf(&order).Elem().FieldByName(field).Interface() == "" && field != "TotalPrice" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Status": 400,
				"error":  stmt,
			})
		}
	}

	//To check if phone number is valid or not
	if len(order.PhoneNumber) != 10 {
		log.Error.Printf("Error : 'Invalid phone number' Status : 400 ")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "Invalid phone number",
		})
	}

	claims := middleware.GetTokenClaims(c)
	UserId, _ := strconv.Atoi(claims["User-id"].(string))
	order.UserId = uint(UserId)
	_, err := repository.ReadProductIdByProductData(db.Connection, order)
	if err != nil {
		log.Error.Printf("Error : 'Product is not found' Status : 404 ")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 404,
			"error":  "Product is not found",
		})
	}
	productPrice, _ := strconv.Atoi(order.ProductPrice)
	ramPrice, _ := strconv.Atoi(order.RamPrice)
	if order.DvdRwDrive {
		order.TotalPrice = strconv.Itoa(productPrice + ramPrice + 3000)
	} else {
		order.TotalPrice = strconv.Itoa(productPrice + ramPrice)
	}
	if err := repository.CreateOrder(db.Connection, order); err != nil {
		log.Error.Printf("Error : '%s' Status : 400\n", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  err.Error(),
		})
	}

	orderId := repository.ReadOrderId(db.Connection)
	var status models.OrderStatus
	status.OrderId = orderId
	status.UserId = order.UserId
	repository.CreateOrderStatus(db.Connection, status)
	URL := fmt.Sprintf("http://:8000/common/getOrderStatus/%v", orderId)
	log.Info.Println("Message : 'Order added successfully' Status : 200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":                           200,
		"message":                          "Order added successfully",
		"click here to get a order status": URL,
	})
}

// Handler for Cancel a order by order-id
func (db Database) CancelOrderById(c echo.Context) error {
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Deleteorder-API called'")
	order, err := repository.ReadOrderByOrderId(db.Connection, c.Param("order_id"))
	if err == nil {
		order.PaymentStatus = "Refunded"
		repository.UpdateOrderById(db.Connection, order)
		repository.DeleteOrderByOrderId(db.Connection, c.Param("order_id"))
		status, _ := repository.ReadOrderStatusByOrderId(db.Connection, order.OrderId)
		status.PaymentStatus = "Refunded"
		status.OrderStatus = "cancelled"
		repository.UpdateOrderStatus(db.Connection, status)
		repository.DeleteOrderStatus(db.Connection, status)
		log.Info.Println("Message : 'order deleted successfully' Status : 200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  200,
			"message": "order deleted Successfully!!!",
		})
	}
	log.Error.Println("Error : 'order not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "order not found",
	})
}

// Handler for get orders
func (db Database) GetOrders(c echo.Context) error {
	log := logs.Log()
	if err := middleware.UserAuth(c); err == nil {
		log.Info.Println("Message : 'GetOrders-API called'")
		claims := middleware.GetTokenClaims(c)
		Orders, err := repository.ReadOrdersByUser(db.Connection, claims["User-id"].(string))
		OrderData := make([]models.OrderProductReq, len(Orders))
		if err == nil && len(Orders) > 0 {
			for index, order := range Orders {
				OrderData[index].BrandName = order.BrandName
				OrderData[index].ProductPrice = order.ProductPrice
				OrderData[index].RamCapacity = order.RamCapacity
				OrderData[index].RamPrice = order.RamPrice
				OrderData[index].DvdRwDrive = order.DvdRwDrive
				OrderData[index].Name = order.Name
				OrderData[index].Address = order.Address
				OrderData[index].PhoneNumber = order.PhoneNumber
				OrderData[index].TotalPrice = order.TotalPrice
			}
			log.Info.Println("Message : 'Order(s) retrieved successfully' Status : 200")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": 200,
				"Orders": OrderData,
			})
		}
		log.Error.Println("Error : 'Order not found' Status : 404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": 404,
			"error":  "Order not found",
		})

	} else if err := middleware.AdminAuth(c); err == nil {
		log.Info.Println("Message : 'GetOrders-API called'")
		Orders, err := repository.ReadOrdersByAdmin(db.Connection)
		OrderData := make([]models.OrderProductReq, len(Orders))
		if err == nil && len(Orders) > 0 {
			for index, order := range Orders {
				OrderData[index].BrandName = order.BrandName
				OrderData[index].ProductPrice = order.ProductPrice
				OrderData[index].RamCapacity = order.RamCapacity
				OrderData[index].RamPrice = order.RamPrice
				OrderData[index].DvdRwDrive = order.DvdRwDrive
				OrderData[index].Name = order.Name
				OrderData[index].Address = order.Address
				OrderData[index].PhoneNumber = order.PhoneNumber
				OrderData[index].TotalPrice = order.TotalPrice
			}
			log.Info.Println("Message : 'Order(s) retrieved successfully' Status : 200")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": 200,
				"Orders": OrderData,
			})
		}
		log.Error.Println("Error : 'Order not found' Status : 404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": 404,
			"error":  "Order not found",
		})
	}
	log.Error.Println("Error : 'unauthorized entry' Status : 401")
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"error":  "unauthorized entry",
		"status": 401,
	})
}

// Payment handler
func (db Database) Payment(c echo.Context) error {
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Payment-API called'")
	order, err := repository.ReadOrderByOrderId(db.Connection, c.Param("order_id"))
	if err != nil {
		log.Error.Println("Error : 'Order not found' Status : 404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": 404,
			"error":  "Order not found",
		})
	}
	if order.PaymentStatus == "pending" {
		var payment models.PaymentReq
		if err := c.Bind(&payment); err != nil {
			log.Error.Println("Error : 'internal server error' Status : 500")
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": 500,
				"error":  "internal server error",
			})
		}
		if payment.Payment == order.TotalPrice {
			order.PaymentStatus = "Paid"
			if err := repository.UpdateOrderById(db.Connection, order); err == nil {
				status, _ := repository.ReadOrderStatusByOrderId(db.Connection, order.OrderId)
				status.PaymentStatus = "paid"
				status.OrderStatus = "order confirmed"
				repository.UpdateOrderStatus(db.Connection, status)
				log.Info.Println("Message : 'Payment successful' Status : 200")
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": 200,
					"Orders": "Payment successful",
				})
			}
		}
		log.Error.Println("Error : 'Payment not matching with the order price' Status : 400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": 400,
			"error":  "Payment not matching with the order price",
		})
	}
	log.Error.Println("Error : 'Already paid' Status : 400")
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": 400,
		"error":  "Already paid",
	})
}

// Handler for update a order status by order-id
func (db Database) UpdateOrderStatusById(c echo.Context) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'UpdateOrderStatus-API called'")
	ord, _ := strconv.Atoi(c.Param("order_id"))
	orderId := uint(ord)
	Status, err := repository.ReadOrderStatusByOrderId(db.Connection, orderId)
	if err == nil {
		if err := c.Bind(&Status); err != nil {
			log.Error.Println("Error : 'internal server error' Status : 500")
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": 500,
				"error":  "internal server error",
			})
		}

		fields := structs.Names(models.OrderStatusReq{})
		for _, field := range fields {
			if reflect.ValueOf(&Status).Elem().FieldByName(field).Interface() == "" {
				log.Error.Println("Error : 'no data found to do update' Status : 404")
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"status": 404,
					"error":  "no data found to do update",
				})
			}
		}

		repository.UpdateOrderStatus(db.Connection, Status)
		log.Info.Println("Message : 'Order status updated successfully' Status : 200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  200,
			"message": "Order status updated Successfully!!!",
		})

	}
	log.Error.Println("Error : 'order not found' Status : 404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": 404,
		"error":  "order not found",
	})
}

// Handler for get order status
func (db Database) GetOrderStatusById(c echo.Context) error {
	log := logs.Log()
	log.Info.Println("Message : 'GetOrderStatus-API called'")
	ord, _ := strconv.Atoi(c.Param("order_id"))
	orderId := uint(ord)
	Status, err := repository.ReadOrderStatusByOrderId(db.Connection, orderId)
	if err != nil {
		log.Error.Println("Error : 'Order not found' Status : 404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": 404,
			"error":  "Order not found",
		})
	}
	order, _ := repository.ReadOrderByOrderId(db.Connection, c.Param("order_id"))
	Status.BrandName = order.BrandName
	Status.Name = order.Name
	Status.Address = order.Address
	Status.PhoneNumber = order.PhoneNumber
	Status.TotalPrice = order.TotalPrice
	if order.DvdRwDrive {
		Status.IncludedProduct = "DVD RW Drive"
	} else {
		Status.IncludedProduct = "None"
	}
	log.Info.Println("Message : 'Order status retrieved successfully' Status : 200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       200,
		"Order Status": Status,
	})
}

// Handler for get all order status
func (db Database) GetAllOrderStatus(c echo.Context) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'GetAllOrderStatus-API called'")
	Statuses, err := repository.ReadOrderStatus(db.Connection)
	if err != nil && len(Statuses) == 0 {
		log.Error.Println("Error : 'Status not found' Status : 404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": 404,
			"error":  "Status not found",
		})
	}
	for index, status := range Statuses {
		orderId := strconv.Itoa(int(status.OrderId))
		order, _ := repository.ReadOrderByOrderId(db.Connection, orderId)
		Statuses[index].BrandName = order.BrandName
		Statuses[index].Name = order.Name
		Statuses[index].Address = order.Address
		Statuses[index].PhoneNumber = order.PhoneNumber
		Statuses[index].TotalPrice = order.TotalPrice
		if order.DvdRwDrive {
			Statuses[index].IncludedProduct = "DVD RW Drive"
		} else {
			Statuses[index].IncludedProduct = "None"
		}
	}
	log.Info.Println("Message : 'Order statuses retrieved successfully' Status : 200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":         200,
		"Order Statuses": Statuses,
	})
}
