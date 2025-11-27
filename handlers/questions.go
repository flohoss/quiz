package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/quiz/config"
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
	Language string `query:"language" doc:"Language code"`
}) (*struct {
	Body config.Quiz `json:"questions"`
}, error) {
	return &struct {
		Body config.Quiz `json:"questions"`
	}{Body: config.GetQuiz(input.Language)}, nil
}

func validateAnswersHandler(ctx context.Context, input *struct {
	Language string              `query:"language" doc:"Language code"`
	Body     []config.QuizAnswer `json:"answers"`
}) (*struct {
	Body config.Quiz `json:"results"`
}, error) {
	if len(input.Body) == 0 {
		return nil, huma.Error400BadRequest("no answers provided")
	}

	results, err := config.ValidateQuizAnswers(input.Body, input.Language)
	if err != nil {
		return nil, huma.Error400BadRequest("validation failed: " + err.Error())
	}

	return &struct {
		Body config.Quiz `json:"results"`
	}{Body: results}, nil
}
