package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VDliveson/SurgeForms/go-producer/constants"
	"github.com/VDliveson/SurgeForms/go-producer/models"
	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func HomeRoute(c *fiber.Ctx) error {
	response := fiber.Map{
		"api":         "Producer API forms route",
		"version":     "1.0",
		"description": "This is the forms route of the Producer API",
	}
	return c.Status(http.StatusOK).JSON(constants.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &response,
	})
}

func CreateForm(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var formCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Form", constants.DatabaseName)
	var qsCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Question", constants.DatabaseName)

	var form constants.FormBody

	err := c.BodyParser(&form)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(constants.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{}})
	}

	validationErr := validate.Struct(&form)
	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(constants.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"details": validationErr.Error()}})
	}

	newForm := models.FormSchema{
		Id:          primitive.NewObjectID(),
		Title:       form.Title,
		Description: form.Description,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	formResult, err := formCollection.InsertOne(ctx, newForm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": err.Error()},
		})
	}

	formId := formResult.InsertedID.(primitive.ObjectID)

	var qsArr []interface{}
	for _, qs := range form.Questions {
		newQuestion := models.QuestionSchema{
			Id:   primitive.NewObjectID(),
			Form: formId,
			Text: qs.Text,
			Type: qs.Type,
		}
		qsArr = append(qsArr, newQuestion)
	}

	qsResult, err := qsCollection.InsertMany(ctx, qsArr)
	if err != nil {
		_, err := formCollection.DeleteOne(ctx, bson.M{"_id": formId})
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": err.Error()},
		})
	}

	var createdForm models.FormSchema
	err = formCollection.FindOne(ctx, bson.M{"_id": formId}).Decode(&createdForm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": "Error retrieving created form details"},
		})
	}

	// Retrieve created questions details
	var createdQuestions []interface{}
	for _, id := range qsResult.InsertedIDs {
		var question models.QuestionSchema
		err := qsCollection.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)}).Decode(&question)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(constants.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"details": "Error retrieving question details"},
			})
		}
		createdQuestions = append(createdQuestions, &fiber.Map{
			"_id": question.Id.Hex(),
			"form": &fiber.Map{
				"_id":         createdForm.Id,
				"title":       createdForm.Title,
				"description": createdForm.Description,
			},
			"text": question.Text,
			"type": question.Type,
		})
	}

	return c.Status(http.StatusOK).JSON(constants.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"createdForm": &fiber.Map{
				"_id":         createdForm.Id,
				"title":       createdForm.Title,
				"description": createdForm.Description,
			},
			"createdQs": createdQuestions,
		},
	})
}

func GetForm(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var formCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Form", constants.DatabaseName)

	formID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(formID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"error": "Invalid form ID format",
			},
		})
	}

	var form models.FormSchema
	err = formCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&form)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(constants.Response{
				Status:  http.StatusNotFound,
				Message: "form not found",
				Data: &fiber.Map{
					"details": err.Error(),
				},
			})
		}
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error fetching form",
			Data: &fiber.Map{
				"details": err.Error(),
			},
		})
	}
	return c.Status(http.StatusOK).JSON(constants.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"_id":         form.Id,
			"title":       form.Title,
			"description": form.Description,
		},
	})
}

func GetQuestion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var formCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Form", constants.DatabaseName)
	var qsCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Question", constants.DatabaseName)

	questionID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"error": "Invalid question ID format",
			},
		})
	}

	var question models.QuestionSchema
	err = qsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(constants.Response{
				Status:  http.StatusNotFound,
				Message: "question not found",
				Data: &fiber.Map{
					"details": err.Error(),
				},
			})
		}
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error fetching question",
			Data: &fiber.Map{
				"details": err.Error(),
			},
		})
	}

	var form models.FormSchema
	err = formCollection.FindOne(ctx, bson.M{"_id": question.Form}).Decode(&form)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(constants.Response{
				Status:  http.StatusNotFound,
				Message: "form not found",
				Data: &fiber.Map{
					"details": err.Error(),
				},
			})
		}
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error fetching form",
			Data: &fiber.Map{
				"details": err.Error(),
			},
		})
	}

	return c.Status(http.StatusOK).JSON(constants.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"question": &fiber.Map{
				"_id": question.Id.Hex(),
				"form": &fiber.Map{
					"_id":         form.Id.Hex(),
					"title":       form.Title,
					"description": form.Description,
				},
				"text": question.Text,
				"type": question.Type,
			},
		},
	})
}

// func CreateResponse(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var formCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Form", constants.DatabaseName)
// 	var qsCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Question", constants.DatabaseName)
// 	var responseCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Response", constants.DatabaseName)

// 	var response constants.Response

// 	err := c.BodyParser(&response)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(constants.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{}})
// 	}
// }
