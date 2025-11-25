package config

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	ConfigFolder = "./config/"
)

var cfg Config

var validate *validator.Validate
var mu sync.RWMutex

type Config struct {
	LogLevel string         `mapstructure:"log_level" validate:"omitempty,oneof=debug info warn error"`
	TimeZone string         `mapstructure:"time_zone" validate:"omitempty,timezone"`
	Server   ServerSettings `mapstructure:"server"`
	Quiz     QuizSettings   `mapstructure:"quiz"`
}

type ServerSettings struct {
	Address string `mapstructure:"address" validate:"required,ipv4"`
	Port    int    `mapstructure:"port" validate:"required,gte=1024,lte=65535"`
}

type QuizSettings struct {
	AmountOfQuestions int        `mapstructure:"amount_of_questions" validate:"required,gte=1"`
	Questions         []Question `mapstructure:"questions" validate:"required,dive"`
}

type Question struct {
	ID            int                 `mapstructure:"id" validate:"min=0"`
	Question      map[string]string   `mapstructure:"question" validate:"required"`
	Answers       map[string][]string `mapstructure:"answers" validate:"required"`
	CorrectAnswer int                 `mapstructure:"correct_answer" validate:"min=0,max=2"`
}

func init() {
	os.Mkdir(ConfigFolder, os.ModePerm)
	validate = validator.New()
}

func New() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("time_zone", "Etc/UTC")
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8156)
	viper.SetDefault("quiz.amount_of_questions", 10)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigFolder)
	viper.SetEnvPrefix("CHRISTMAS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.WriteConfigAs(ConfigFolder + "config.yaml")
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		} else {
			slog.Error("Failed to read configuration file", "error", err)
			os.Exit(1)
		}
	}

	viper.AutomaticEnv()

	if err := ValidateAndLoadConfig(viper.GetViper()); err != nil {
		slog.Error("Initial configuration validation failed", "error", err)
		os.Exit(1)
	}
}

func ValidateAndLoadConfig(v *viper.Viper) error {
	var tempCfg Config
	if err := v.Unmarshal(&tempCfg); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := validate.Struct(tempCfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	mu.Lock()
	cfg = tempCfg
	mu.Unlock()

	os.Setenv("TZ", cfg.TimeZone)
	return nil
}

func ConfigLoaded() bool {
	return viper.ConfigFileUsed() != ""
}

func GetLogLevel() slog.Level {
	mu.RLock()
	defer mu.RUnlock()
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func GetServer() string {
	mu.RLock()
	defer mu.RUnlock()
	return fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
}

func validateLanguage(lang string) string {
	if lang == "en" || lang == "de" {
		return lang
	}
	return "en"
}

type QuestionAndAnswer struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

func GetQuiz(lang string) []QuestionAndAnswer {
	lang = validateLanguage(lang)

	mu.RLock()
	defer mu.RUnlock()

	amount := min(cfg.Quiz.AmountOfQuestions, len(cfg.Quiz.Questions))

	shuffled := shuffleQuestions(cfg.Quiz.Questions, amount)

	result := make([]QuestionAndAnswer, len(shuffled))
	for i, q := range shuffled {
		question := q.Question[lang]
		answers := q.Answers[lang]

		result[i] = QuestionAndAnswer{
			ID:       q.ID,
			Question: question,
			Answers:  answers,
		}
	}

	return result
}

type ValidationResult struct {
	ID            int     `json:"id"`
	Question      string  `json:"question"`
	UserAnswer    string  `json:"user_answer"`
	CorrectAnswer *string `json:"correct_answer"`
}

type QuizAnswer struct {
	ID     int `json:"id"`
	Answer int `json:"answer"`
}

func ValidateQuizAnswers(answers []QuizAnswer, lang string) []ValidationResult {
	lang = validateLanguage(lang)

	mu.RLock()
	defer mu.RUnlock()

	results := make([]ValidationResult, len(answers))

	for i, answer := range answers {
		result := ValidationResult{
			ID: answer.ID,
		}

		var foundQuestion *Question
		for j := range cfg.Quiz.Questions {
			if cfg.Quiz.Questions[j].ID == answer.ID {
				foundQuestion = &cfg.Quiz.Questions[j]
				break
			}
		}

		if foundQuestion != nil {
			question := foundQuestion.Question[lang]
			answers := foundQuestion.Answers[lang]

			result.Question = question

			// Set user answer as string
			if answer.Answer >= 0 && answer.Answer < len(answers) {
				result.UserAnswer = answers[answer.Answer]
			} else {
				result.UserAnswer = "Invalid answer"
			}

			// If answer is incorrect, set correct_answer to the correct answer string
			// If answer is correct, correct_answer remains nil
			if answer.Answer != foundQuestion.CorrectAnswer {
				if foundQuestion.CorrectAnswer >= 0 && foundQuestion.CorrectAnswer < len(answers) {
					correctAnswerStr := answers[foundQuestion.CorrectAnswer]
					result.CorrectAnswer = &correctAnswerStr
				}
			}
		} else {
			result.Question = "Question not found"
			result.UserAnswer = "Invalid question"
			notFoundStr := "Question not found"
			result.CorrectAnswer = &notFoundStr
		}

		results[i] = result
	}

	return results
}

func shuffleQuestions(questions []Question, amount int) []Question {
	shuffled := make([]Question, len(questions))
	copy(shuffled, questions)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	if amount > len(shuffled) {
		amount = len(shuffled)
	}

	return shuffled[:amount]
}
