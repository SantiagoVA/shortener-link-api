package validations

import (
	"context"
	"shortener-app/database"
	"shortener-app/models"

	"github.com/alexedwards/argon2id"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = database.GetCollection("users")

func SamePassword(password, email string) (bool, *models.User) {
	var structure *models.User
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "email", Value: email}}).Decode(&structure)

	if err != nil {
		return false, &models.User{}
	}

	match, _ := argon2id.ComparePasswordAndHash(password, structure.Password)

	return match, structure
}

func IsUsedBefore(email string) bool {
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "email", Value: email}})
	return err == nil
}

func UsedBeforeUpdate(email, id string) bool {
	var user *models.User
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		return false
	}

	if user.Id != id && user.Email == email {
		return true
	}

	return false
}
