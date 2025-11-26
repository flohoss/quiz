package config

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"slices"
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
	LogLevel  string         `mapstructure:"log_level" validate:"omitempty,oneof=debug info warn error"`
	TimeZone  string         `mapstructure:"time_zone" validate:"omitempty,timezone"`
	Server    ServerSettings `mapstructure:"server"`
	UI        UI             `mapstructure:"ui"`
	Quiz      QuizSettings   `mapstructure:"quiz"`
	Languages []string       `mapstructure:"-" validate:"-"`
}

type ServerSettings struct {
	Address string `mapstructure:"address" validate:"required,ipv4"`
	Port    int    `mapstructure:"port" validate:"required,gte=1024,lte=65535"`
}

type UI struct {
	Title string `mapstructure:"title"`
	Logo  string `mapstructure:"logo"`
}

type QuizSettings struct {
	AmountOfQuestions int        `mapstructure:"amount_of_questions" validate:"required,gte=1"`
	Questions         []Question `mapstructure:"questions" validate:"required,dive"`
}

type Question struct {
	ID            int                 `mapstructure:"id" validate:"min=0"`
	Question      map[string]string   `mapstructure:"question" validate:"required"`
	Answers       map[string][]string `mapstructure:"answers" validate:"required"`
	CorrectAnswer int                 `mapstructure:"correct_answer" validate:"min=1,max=3"`
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

	langSet := make(map[string]struct{})
	for _, q := range tempCfg.Quiz.Questions {
		for lang := range q.Question {
			langSet[lang] = struct{}{}
		}
	}
	tempCfg.Languages = make([]string, 0, len(langSet))
	for lang := range langSet {
		tempCfg.Languages = append(tempCfg.Languages, lang)
	}

	for _, q := range tempCfg.Quiz.Questions {
		for _, lang := range tempCfg.Languages {
			if _, ok := q.Question[lang]; !ok {
				return fmt.Errorf("question id %d is missing language key: %s", q.ID, lang)
			}
		}
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

func GetUI() UI {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.UI
}

func ValidateLanguage(lang string) error {
	mu.RLock()
	defer mu.RUnlock()
	if slices.Contains(cfg.Languages, lang) {
		return nil
	}
	return fmt.Errorf("unsupported language, supported languages are: %v", cfg.Languages)
}

type QuestionAndAnswer struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers" nullable:"false"`
	Answer   *int     `json:"answer,omitempty"`
	Correct  *bool    `json:"correct,omitempty"`
}

func GetQuiz(lang string) Quiz {
	mu.RLock()
	defer mu.RUnlock()

	amount := min(cfg.Quiz.AmountOfQuestions, len(cfg.Quiz.Questions))
	shuffled := shuffleQuestions(cfg.Quiz.Questions, amount)

	quiz := Quiz{
		Questions: make([]QuestionAndAnswer, 0, amount),
		Total:     amount,
	}

	for _, q := range shuffled {
		question := q.Question[lang]
		answers := q.Answers[lang]

		quiz.Questions = append(quiz.Questions, QuestionAndAnswer{
			ID:       q.ID,
			Question: question,
			Answers:  answers,
		})
	}

	return quiz
}

type Quiz struct {
	Questions []QuestionAndAnswer `json:"questions" nullable:"false"`
	Correct   *int                `json:"correct,omitempty"`
	Total     int                 `json:"total"`
}

type QuizAnswer struct {
	ID     int `json:"id"`
	Answer int `json:"answer"`
}

func ValidateQuizAnswers(answers []QuizAnswer, lang string) (Quiz, error) {
	mu.RLock()
	defer mu.RUnlock()

	quiz := Quiz{
		Questions: make([]QuestionAndAnswer, 0, len(answers)),
		Total:     len(answers),
		Correct:   new(int),
	}

	for _, answer := range answers {
		var foundQuestion *Question
		for j := range cfg.Quiz.Questions {
			if cfg.Quiz.Questions[j].ID == answer.ID {
				foundQuestion = &cfg.Quiz.Questions[j]
				break
			}
		}

		if foundQuestion == nil {
			continue
		}

		answers := foundQuestion.Answers[lang]

		if answer.Answer < 1 || answer.Answer > len(answers) {
			return quiz, fmt.Errorf("invalid answer index %d for question %d", answer.Answer, answer.ID)
		}

		question := QuestionAndAnswer{
			ID:       answer.ID,
			Answer:   &answer.Answer,
			Question: foundQuestion.Question[lang],
			Answers:  answers,
			Correct:  new(bool),
		}

		if answer.Answer == foundQuestion.CorrectAnswer {
			*quiz.Correct++
			*question.Correct = true
		}

		quiz.Questions = append(quiz.Questions, question)
	}

	return quiz, nil
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
