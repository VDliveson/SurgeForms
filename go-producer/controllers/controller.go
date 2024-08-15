package controllers

import (
	"context"
	"encoding/json"
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

func CreateResponse(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var formCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Form", constants.DatabaseName)
	var qsCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Question", constants.DatabaseName)
	var responseCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Response", constants.DatabaseName)
	var ansCollection *mongo.Collection = utils.GetCollection(utils.DBClient, "Answer", constants.DatabaseName)

	var service string = c.Get("service")
	var responseBody constants.ResponseBody

	err := c.BodyParser(&responseBody)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(constants.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{},
		})
	}

	formID, err := primitive.ObjectIDFromHex(responseBody.Form)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(constants.Response{
			Status:  http.StatusNotFound,
			Message: "Internal Server Error",
			Data:    &fiber.Map{},
		})
	}

	userID, err := primitive.ObjectIDFromHex(responseBody.User)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(constants.Response{
			Status:  http.StatusNotFound,
			Message: "Internal Server Error",
			Data:    &fiber.Map{},
		})
	}

	var form models.FormSchema
	err = formCollection.FindOne(ctx, bson.M{"_id": formID}).Decode(&form)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(constants.Response{
			Status:  http.StatusNotFound,
			Message: "Form not found",
			Data:    &fiber.Map{},
		})
	}

	responseSchema := models.ResponseSchema{
		Id:          primitive.NewObjectID(),
		Form:        formID,
		User:        userID,
		SubmittedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	responseResult, err := responseCollection.InsertOne(ctx, responseSchema)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": err.Error()},
		})
	}

	responseId := responseResult.InsertedID.(primitive.ObjectID)

	var ansArr []interface{}
	var qsArr []interface{}
	for _, answer := range responseBody.Answers {
		var question models.QuestionSchema
		questionId, _ := primitive.ObjectIDFromHex(answer.Question)
		err = qsCollection.FindOne(ctx, bson.M{"_id": questionId}).Decode(&question)
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

		newAnswer := models.AnswerSchema{
			Id:       primitive.NewObjectID(),
			Question: questionId,
			Response: responseId,
			Text:     answer.Text,
		}
		ansArr = append(ansArr, newAnswer)
		qsArr = append(qsArr, question)
	}

	ansResult, err := ansCollection.InsertMany(ctx, ansArr)
	if err != nil {
		_, err := responseCollection.DeleteOne(ctx, bson.M{"_id": responseId})
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": err.Error()},
		})
	}

	var createdAnswers []interface{}

	for idx, val := range ansResult.InsertedIDs {
		question := qsArr[idx].(models.QuestionSchema)

		var answer models.AnswerSchema
		err = ansCollection.FindOne(ctx, bson.M{"_id": val}).Decode(&answer)
		if err != nil {
			_, err := responseCollection.DeleteOne(ctx, bson.M{"_id": responseId})
			return c.Status(http.StatusInternalServerError).JSON(constants.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"details": err.Error()},
			})
		}

		createdAnswers = append(createdAnswers, &fiber.Map{
			"_id": val,
			"question": &fiber.Map{
				"_id":  question.Id,
				"form": question.Form,
				"text": question.Text,
				"type": question.Type,
			},
			"response": answer.Response,
			"text":     answer.Text,
		})
	}

	formData := &fiber.Map{
		"_id":   form.Id.Hex(),
		"title": form.Title,
	}

	createdResponse := &fiber.Map{
		"_id":  responseId,
		"form": formData,
		"user": responseBody.User,
	}

	data := &fiber.Map{
		"createdResponse": createdResponse,
		"createdAnswers":  createdAnswers,
		"metadata":        responseBody.Metadata,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		_, err := responseCollection.DeleteOne(ctx, bson.M{"_id": responseId})
		return c.Status(http.StatusInternalServerError).JSON(constants.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"details": err.Error()},
		})
	}

	utils.SendData(jsonData, service) // Call rabbitmq service

	return c.Status(http.StatusOK).JSON(constants.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"createdResponse": (*data)["createdResponse"],
			"createdAnswers":  (*data)["createdAnswers"],
		},
	})
}
