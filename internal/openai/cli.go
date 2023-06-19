package openai

import (
	"errors"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/christian-gama/autocommit/internal/helpers"
)

// AskConfigsCli is a command line interface that asks the user for the configuration.
type AskConfigsCli interface {
	Execute() (*Config, error)
}

// askConfigsCliImpl is the implementation of AskConfigsCli.
type askConfigsCliImpl struct{}

// Execute asks the user for the configuration.
func (a *askConfigsCliImpl) Execute() (*Config, error) {
	questions := helpers.CreateQuestions(
		a.createApiKeyQuestion,
		a.createModelQuestion,
		a.createTemperatureQuestion,
	)

	type Answers struct {
		Model        string
		OpenAIAPIKey string
		Temperature  float32
	}

	var answers Answers

	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	return NewConfig(answers.OpenAIAPIKey, answers.Model, answers.Temperature), nil
}

func (a *askConfigsCliImpl) createModelQuestion() *survey.Question {
	prompt := survey.Select{
		Message: "Model name",
		Help:    "A model can be an algorithm or a set of algorithms that have been trained on data to make predictions or decisions.",
		Default: GPT3Dot5Turbo16k,
		Options: AllowedModels,
		Description: func(value string, index int) string {
			if value == GPT4 || value == GPT432K {
				return "Beta - May not be available for all users"
			}
			return ""
		},

		VimMode: true,
	}

	return &survey.Question{
		Name:   "Model",
		Prompt: &prompt,
	}
}

func (a *askConfigsCliImpl) createApiKeyQuestion() *survey.Question {
	prompt := survey.Password{
		Message: "OpenAI API Key",
		Help:    "The OpenAI API Key is used to authenticate your requests to the OpenAI API.",
	}

	return &survey.Question{
		Name:   "OpenAIAPIKey",
		Prompt: &prompt,
		Validate: func(ans interface{}) error {
			return ValidateApiKey(ans.(string))
		},
	}
}

func (a *askConfigsCliImpl) createTemperatureQuestion() *survey.Question {
	prompt := survey.Input{
		Message: "Temperature",
		Help:    "Temperature refers to a parameter that controls the randomness of the output generated by the model.",
		Default: "0.3",
	}

	return &survey.Question{
		Name:      "Temperature",
		Prompt:    &prompt,
		Transform: a.transformToFloat32,
		Validate: func(ans interface{}) error {
			f, err := strconv.ParseFloat(ans.(string), 32)
			if err != nil {
				return errors.New("Invalid temperature - must be a number")
			}

			return ValidateTemperature(float32(f))
		},
	}
}

func (a *askConfigsCliImpl) transformToFloat32(ans interface{}) (newAns interface{}) {
	f, err := strconv.ParseFloat(ans.(string), 32)
	if err != nil {
		return errors.New("Invalid temperature - must be a number")
	}

	return float32(f)
}

// NewAskConfigsCli creates a new instance of AskConfigsCli.
func NewAskConfigsCli() AskConfigsCli {
	return &askConfigsCliImpl{}
}

type AskToChangeModelCli interface {
	Execute() (bool, error)
}

type askToChangeModelCliImpl struct{}

func (a *askToChangeModelCliImpl) Execute() (bool, error) {
	questions := helpers.CreateQuestions(
		a.createModelQuestion,
	)

	type Answers struct {
		ChangeModel bool
	}

	var answers Answers

	err := survey.Ask(questions, &answers)
	if err != nil {
		return false, err
	}

	return answers.ChangeModel, nil
}

func (a *askToChangeModelCliImpl) createModelQuestion() *survey.Question {
	prompt := survey.Confirm{
		Message: "You reached the maximum number of tokens, but there is a model that can generate longer messages. Do you want to temporarily change the model?",
		Help:    "A model have a limited amount of tokens that can be generated at once. If you want to generate longer messages, you can temporarily change the model.",
		Default: true,
	}

	return &survey.Question{
		Name:   "ChangeModel",
		Prompt: &prompt,
	}
}

func NewAskToChangeModelCli() AskToChangeModelCli {
	return &askToChangeModelCliImpl{}
}
