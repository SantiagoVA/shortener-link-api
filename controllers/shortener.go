package controllers

import (
	"context"
	"fmt"
	"shortener-app/database"
	"shortener-app/functions"
	"shortener-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	linksCollection = database.GetCollection("links")
)

func Default(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "Hello world!",
		"data":         nil,
		"token":        nil,
		"refreshToken": nil,
	})
}

func NewShorted(c *fiber.Ctx) error {
	link := new(models.Links)
	if err := c.BodyParser(link); err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	if !functions.IsUrl(link.InputLink) {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Bad request. The link entered is invalid",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	link.Id, _ = shortid.Generate()
	link.OutputLink = c.BaseURL() + "/l/" + link.Id

	_, err := linksCollection.InsertOne(context.Background(), link)

	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Internal server error",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}

	token := c.Get("token", "")
	if token != "" {
		isValid, claims := functions.IsValidToken(token)
		if !isValid {
			c.Status(401)
			return c.JSON(map[string]interface{}{
				"error":        true,
				"message":      "Invalid token",
				"data":         nil,
				"token":        nil,
				"refreshToken": nil,
			})
		}
		var user *models.User
		id := claims["id"].(string)
		err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "id", Value: id}}).Decode(&user)

		if err != nil {
			c.Status(500)
			return c.JSON(map[string]interface{}{
				"error":        true,
				"message":      "Internal server error",
				"data":         nil,
				"token":        nil,
				"refreshToken": nil,
			})
		}

		user.Links = append(user.Links, link)

		filter := bson.M{"id": id}

		update := bson.M{
			"$set": bson.M{
				"links": user.Links,
			},
		}

		_, err = userCollection.UpdateOne(context.Background(), filter, update)

		if err != nil {
			c.Status(500)
			return c.JSON(map[string]interface{}{
				"error":        true,
				"message":      "Internal server error",
				"data":         nil,
				"token":        nil,
				"refreshToken": nil,
			})
		}
	}

	c.Status(201)
	return c.JSON(map[string]interface{}{
		"error":        false,
		"message":      "Link created",
		"data":         link,
		"token":        nil,
		"refreshToken": nil,
	})
}

func Redirect(c *fiber.Ctx) error {
	id := c.Params("id")
	var link *models.Links
	err := linksCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&link)
	if err != nil {
		c.Status(404)
		return c.JSON(map[string]interface{}{
			"error":        true,
			"message":      "Link not found",
			"data":         nil,
			"token":        nil,
			"refreshToken": nil,
		})
	}
	return c.Redirect(link.InputLink, 301)
}
