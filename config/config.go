package config

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"regexp"
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
	LogLevel  string            `mapstructure:"log_level" yaml:"log_level" validate:"omitempty,oneof=debug info warn error"`
	TimeZone  string            `mapstructure:"time_zone" yaml:"time_zone" validate:"omitempty,timezone"`
	Server    ServerSettings    `mapstructure:"server" yaml:"server"`
	App       AppSettings       `mapstructure:"app" yaml:"app"`
	Questions []QuestionSetting `mapstructure:"questions" yaml:"questions" validate:"dive"`
}

type ServerSettings struct {
	Address string `mapstructure:"address" yaml:"address" validate:"required,ipv4"`
	Port    int    `mapstructure:"port" yaml:"port" validate:"required,gte=1024,lte=65535"`
}

type AppSettings struct {
	Title             string            `json:"title" mapstructure:"title" yaml:"title" validate:"required"`
	AmountOfQuestions int               `json:"amount_of_questions" mapstructure:"amount_of_questions" yaml:"amount_of_questions" validate:"required,gte=1"`
	Languages         []string          `json:"languages" mapstructure:"languages" yaml:"languages" validate:"dive,required,bcp47_language_tag" nullable:"false"`
	CSSVariables      map[string]string `json:"css_variables" mapstructure:"css_variables" yaml:"css_variables" validate:"dive,required,hexcolor" nullable:"false"`
	Logo              string            `json:"logo" mapstructure:"logo" yaml:"logo" validate:"required,image" `
	Favicon           string            `json:"favicon" mapstructure:"favicon" yaml:"favicon" validate:"required,image"`
	Icons             map[string]string `json:"icons" mapstructure:"icons" yaml:"icons" validate:"required,dive,required,svg" nullable:"false"`
}

type QuestionSetting struct {
	ID            int                 `mapstructure:"id" yaml:"id" validate:"min=0"`
	Question      map[string]string   `mapstructure:"question" yaml:"question" validate:"dive,required"`
	Answers       map[string][]string `mapstructure:"answers" yaml:"answers" validate:"dive,dive,required"`
	CorrectAnswer int                 `mapstructure:"correct_answer" yaml:"correct_answer" validate:"min=1,max=3"`
}

type Quiz struct {
	Questions []QuestionAndAnswer `json:"questions" nullable:"false"`
	Correct   *int                `json:"correct,omitempty"`
	Wrong     *int                `json:"wrong,omitempty"`
	Total     int                 `json:"total"`
}

type QuestionAndAnswer struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers" nullable:"false"`
	Answer   *int     `json:"answer,omitempty"`
	Correct  int      `json:"correct"`
}

func init() {
	os.Mkdir(ConfigFolder, os.ModePerm)
	validate = validator.New()
	validate.RegisterValidation("svg", isSVGString)
}

func New() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("time_zone", "Europe/Berlin")

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8156)

	viper.SetDefault("app.title", "Quiz")
	viper.SetDefault("app.amount_of_questions", 5)
	viper.SetDefault("app.languages", []string{"en", "de"})
	viper.SetDefault("app.css_variables", map[string]string{
		"--color-primary":           "#294221",
		"--color-primary-content":   "#ffffff",
		"--color-secondary":         "#ac3e31",
		"--color-secondary-content": "#ffffff",
		"--color-success":           "#294221",
		"--color-success-content":   "#ffffff",
		"--color-error":             "#ac3e31",
		"--color-error-content":     "#ffffff",
	})
	viper.SetDefault("app.logo", "/app/config/logo.svg")
	viper.SetDefault("app.favicon", "/app/config/logo.svg")
	viper.SetDefault("app.icons", map[string]string{
		"next":     `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M480 96c17.7 0 32 14.3 32 32s-14.3 32-32 32H272V96zM320 288c17.7 0 32 14.3 32 32s-14.3 32-32 32h-64c-17.7 0-32-14.3-32-32s14.3-32 32-32zm64-64c0 17.7-14.3 32-32 32h-48c-17.7 0-32-14.3-32-32s14.3-32 32-32h48c17.7 0 32 14.3 32 32m-96 160c17.7 0 32 14.3 32 32s-14.3 32-32 32h-64c-17.7 0-32-14.3-32-32s14.3-32 32-32zm-88-96h.6c-5.4 9.4-8.6 20.3-8.6 32c0 13.2 4 25.4 10.8 35.6c-24.9 8.7-42.8 32.5-42.8 60.4c0 11.7 3.1 22.6 8.6 32H160C71.6 448 0 376.4 0 288v-61.7c0-42.4 16.9-83.1 46.9-113.1l11.6-11.6C82.5 77.5 115.1 64 149 64h27c35.3 0 64 28.7 64 64v88c0 22.1-17.9 40-40 40s-40-17.9-40-40v-56c0-8.8-7.2-16-16-16s-16 7.2-16 16v56c0 39.8 32.2 72 72 72"/></svg>`,
		"previous": `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M32 96c-17.7 0-32 14.3-32 32s14.3 32 32 32h208V96zm160 192c-17.7 0-32 14.3-32 32s14.3 32 32 32h64c17.7 0 32-14.3 32-32s-14.3-32-32-32zm-64-64c0 17.7 14.3 32 32 32h48c17.7 0 32-14.3 32-32s-14.3-32-32-32h-48c-17.7 0-32 14.3-32 32m96 160c-17.7 0-32 14.3-32 32s14.3 32 32 32h64c17.7 0 32-14.3 32-32s-14.3-32-32-32zm88-96h-.6c5.4 9.4 8.6 20.3 8.6 32c0 13.2-4 25.4-10.8 35.6c24.9 8.7 42.8 32.5 42.8 60.4c0 11.7-3.1 22.6-8.6 32h8.6c88.4 0 160-71.6 160-160v-61.7c0-42.4-16.9-83.1-46.9-113.1l-11.6-11.6C429.5 77.5 396.9 64 363 64h-27c-35.3 0-64 28.7-64 64v88c0 22.1 17.9 40 40 40s40-17.9 40-40v-56c0-8.8 7.2-16 16-16s16 7.2 16 16v56c0 39.8-32.2 72-72 72"/></svg>`,
		"submit":   `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 512"><path fill="currentColor" d="M32 32C14.3 32 0 46.3 0 64s14.3 32 32 32v160c0 53 43 96 96 96v32h64v-32h192v32h64v-32c53 0 96-43 96-96v-96c17.7 0 32-14.3 32-32s-14.3-32-32-32h-64c-17.7 0-32 14.3-32 32v41.3c0 30.2-24.5 54.7-54.7 54.7c-75.5 0-145.6-38.9-185.6-102.9l-4.3-6.9C174.2 67.6 125 37.6 70.7 32.7c-2.2-.5-4.4-.7-6.7-.7zm608 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v8c0 13.3-10.7 24-24 24H64c-17.7 0-32 14.3-32 32s14.3 32 32 32h488c48.6 0 88-39.4 88-88z"/></svg>`,
		"restart":  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M0 224c0 17.7 14.3 32 32 32s32-14.3 32-32c0-53 43-96 96-96h160v32c0 12.9 7.8 24.6 19.8 29.6s25.7 2.2 34.9-6.9l64-64c12.5-12.5 12.5-32.8 0-45.3l-64-64c-9.2-9.2-22.9-11.9-34.9-6.9S320 19.1 320 32v32H160C71.6 64 0 135.6 0 224m512 64c0-17.7-14.3-32-32-32s-32 14.3-32 32c0 53-43 96-96 96H192v-32c0-12.9-7.8-24.6-19.8-29.6s-25.7-2.2-34.9 6.9l-64 64c-12.5 12.5-12.5 32.8 0 45.3l64 64c9.2 9.2 22.9 11.9 34.9 6.9s19.8-16.6 19.8-29.6v-32h160c88.4 0 160-71.6 160-160z"/></svg>`,
	})

	viper.SetDefault("questions", []QuestionSetting{
		{
			ID: 1,
			Question: map[string]string{
				"en": "What is the capital of France?",
				"de": "Was ist die Hauptstadt von Frankreich?",
			},
			Answers: map[string][]string{
				"en": {
					"Berlin",
					"Madrid",
					"Paris",
				},
				"de": {
					"Berlin",
					"Madrid",
					"Paris",
				},
			},
			CorrectAnswer: 3,
		},
	})

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

	for _, q := range tempCfg.Questions {
		for _, lang := range tempCfg.App.Languages {
			if _, ok := q.Question[lang]; !ok {
				return fmt.Errorf("question id %d is missing question language key: %s", q.ID, lang)
			}
			if _, ok := q.Answers[lang]; !ok {
				return fmt.Errorf("question id %d is missing answers language key: %s", q.ID, lang)
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

func GetApp() AppSettings {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.App
}

func ValidateLanguage(lang string) error {
	mu.RLock()
	defer mu.RUnlock()
	if slices.Contains(cfg.App.Languages, lang) {
		return nil
	}
	return fmt.Errorf("unsupported language, supported languages are: %v", cfg.App.Languages)
}

func GetQuiz(lang string) Quiz {
	mu.RLock()
	defer mu.RUnlock()

	amount := min(cfg.App.AmountOfQuestions, len(cfg.Questions))
	shuffled := shuffleQuestions(cfg.Questions, amount)

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
		Wrong:     new(int),
	}

	for _, answer := range answers {
		var foundQuestion *QuestionSetting
		for j := range cfg.Questions {
			if cfg.Questions[j].ID == answer.ID {
				foundQuestion = &cfg.Questions[j]
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

		if answer.Answer == foundQuestion.CorrectAnswer {
			*quiz.Correct++
		} else {
			*quiz.Wrong++
		}

		question := QuestionAndAnswer{
			ID:       answer.ID,
			Answer:   &answer.Answer,
			Question: foundQuestion.Question[lang],
			Answers:  answers,
			Correct:  foundQuestion.CorrectAnswer,
		}

		quiz.Questions = append(quiz.Questions, question)
	}

	return quiz, nil
}

func shuffleQuestions(questions []QuestionSetting, amount int) []QuestionSetting {
	shuffled := make([]QuestionSetting, len(questions))
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

func isSVGString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`(?s)^\s*<svg[\s\S]*?</svg>\s*$`, value)
	return matched
}
