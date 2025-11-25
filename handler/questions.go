package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/christmas/config"
)

func getQuestionsOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-questions",
		Method:      http.MethodGet,
		Path:        "/api/questions",
		Summary:     "Get quiz questions",
		Description: "Get random quiz questions with language support.",
		Tags:        []string{"Quiz"},
	}
}

func validateAnswersOperation() huma.Operation {
	return huma.Operation{
		OperationID: "validate-answers",
		Method:      http.MethodPost,
		Path:        "/api/questions",
		Summary:     "Validate quiz answers",
		Description: "Validate quiz answers and return results with correct answers.",
		Tags:        []string{"Quiz"},
	}
}

func getQuestionsHandler(ctx context.Context, input *struct {
	Lang string `query:"lang" default:"en" enum:"en,de" doc:"Language code"`
}) (*struct {
	Body []config.QuestionAndAnswer `json:"questions"`
}, error) {
	questions, err := config.GetQuiz(input.Lang)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &struct {
		Body []config.QuestionAndAnswer `json:"questions"`
	}{Body: questions}, nil
}

func validateAnswersHandler(ctx context.Context, input *struct {
	Lang string              `query:"lang" default:"en" enum:"en,de" doc:"Language code"`
	Body []config.QuizAnswer `json:"answers"`
}) (*struct {
	Body []config.ValidationResult `json:"results"`
}, error) {
	if len(input.Body) == 0 {
		return nil, huma.Error400BadRequest("No answers provided")
	}

	results, err := config.ValidateQuizAnswers(input.Body, input.Lang)
	if err != nil {
		return nil, huma.Error400BadRequest("Validation failed: " + err.Error())
	}
	return &struct {
		Body []config.ValidationResult `json:"results"`
	}{Body: results}, nil
}
