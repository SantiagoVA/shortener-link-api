package controllers

import (
	"context"
	"shortener-app/database"
	"shortener-app/functions"
	"shortener-app/functions/validations"
	"shortener-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
)

var userCollection = database.GetCollection("users")

func SignUp(c *fiber.Ctx) error {
	token := c.Get("token", "")
	if token != "" {
		isValid, claims := functions.IsValidToken(token)
		if isValid {
			user := models.User{
				Id:    claims["id"].(string),
				Email: claims["email"].(string),
				Name:  claims["name"].(string),
			}
			newToken, err := functions.NewToken(&user)
			if err != nil {
				c.Status(500)
				return c.JSON(map[string]interface{}{
					"error":        true,
					"message":      "The user is already logged in, but a problem happen making the refresh token function",
					"data":         nil,
					"token":        nil,
					"refreshToken": nil,
				})
			}
			c.Status(403)
			return c.JSON(map[string]interface{}{
				"error":        true,
				"message":      "You are already logged in",
				"data":         newToken,
				"token":        newToken,
				"refreshToken": newToken,
			})
		}
	}
	user := new(models.User)
	err := c.BodyParser(user)
	if err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	if !functions.IsEmail(user.Email) {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request. The email is invalid",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	if validations.IsUsedBefore(user.Email) {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "The email is already used",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	user.Id, _ = shortid.Generate()
	user.Password, err = functions.Encrypt(user.Password)

	if err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Error encrypting the password",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}
	_, err = userCollection.InsertOne(context.Background(), user)

	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Internal server error. Error saving the data of the user",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	token, err = functions.NewToken(user)

	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Error generating the token",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	c.Status(201)
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "User created",
		"data":         token,
		"token":        token,
		"refreshToken": token,
	})
}

func Login(c *fiber.Ctx) error {
	token := c.Get("token", "")
	if token != "" {
		isValid, claims := functions.IsValidToken(token)
		if isValid {
			user := models.User{
				Id:    claims["id"].(string),
				Email: claims["email"].(string),
				Name:  claims["name"].(string),
			}
			newToken, err := functions.NewToken(&user)
			if err != nil {
				c.Status(500)
				return c.JSON(map[string]interface{}{
					"error":        true,
					"message":      "The user is already logged in, but a problem happen making the refresh token function",
					"data":         nil,
					"token":        nil,
					"refreshToken": nil,
				})
			}
			c.Status(403)
			return c.JSON(map[string]interface{}{
				"error":        true,
				"message":      "You are already logged in",
				"data":         newToken,
				"token":        newToken,
				"refreshToken": newToken,
			})
		}
	}
	user := new(models.Login)
	err := c.BodyParser(user)
	if err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	isEqual, allUserData := validations.SamePassword(user.Password, user.Email)

	if !isEqual {
		c.Status(401)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad credentials",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	token, err = functions.NewToken(allUserData)

	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Error generating the token",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	c.Status(200)
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "User logged in",
		"data":         token,
		"token":        token,
		"refreshToken": token,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	toUpdate := new(models.User)
	err := c.BodyParser(toUpdate)
	if err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	emailIsUsed := validations.UsedBeforeUpdate(toUpdate.Email, toUpdate.Id)

	token := c.Get("token", "")
	_, claims := functions.IsValidToken(token)
	user := models.User{
		Id:    claims["id"].(string),
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
	}
	newToken, err := functions.NewToken(&user)

	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Error generating the token",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	if emailIsUsed {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request. The email is already used",
			"data":         newToken,
			"token":        newToken,
			"refreshToken": newToken,
		})
	}

	update := bson.M{"$set": bson.M{
		"email":    toUpdate.Email,
		"name":     toUpdate.Name,
		"password": toUpdate.Password,
	}}
	_, err = userCollection.UpdateOne(context.Background(), bson.M{"id": claims["id"].(string)}, update)
	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Internal server error. Error updating the user",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	c.Status(200)
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "User updated",
		"data":         newToken,
		"token":        newToken,
		"refreshToken": newToken,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	token := c.Get("token", "")
	_, claims := functions.IsValidToken(token)
	_, err := userCollection.DeleteOne(context.Background(), bson.M{"id": claims["id"].(string)})

	if err != nil {
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Internal server error. Error deleting the user",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	c.Status(200)
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "User deleted",
		"data":         nil,
		"token":        nil,
		"refreshToken": nil,
	})
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("token", "")
	valid, _ := functions.IsValidToken(token)
	if token == "" || !valid {
		c.Status(401)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Unauthorized. You need to be logged in",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}
	return c.Next()
}

func ListLinks(c *fiber.Ctx) error {
	token := c.Get("token", "")
	_, claims := functions.IsValidToken(token)
	id := claims["id"].(string)
	var user *models.User
	err := linksCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		c.Status(404)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "User not found",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "Links found",
		"data":         user.Links,
		"token":        nil,
		"refreshToken": nil,
	})
}
