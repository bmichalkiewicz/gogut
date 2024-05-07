package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/bmichalkiewicz/gogut/config"
	"github.com/bmichalkiewicz/gogut/facts"

	"github.com/sashabaranov/go-openai"
)

const noexec = "[noexec]"

type Engine struct {
	mode         EngineMode
	config       *config.Config
	client       *openai.Client
	execMessages []openai.ChatCompletionMessage
	chatMessages []openai.ChatCompletionMessage
	channel      chan EngineChatStreamOutput
	pipe         string
	running      bool
}

func NewEngine(mode EngineMode, config *config.Config) (*Engine, error) {
	var client *openai.Client

	if config.GetAIConfig().GetURL() == "" {
		client = openai.NewClient(config.GetAIConfig().GetKey())
	} else {
		clientConfig := openai.DefaultConfig(config.GetAIConfig().GetKey())

		url, err := url.Parse(config.GetAIConfig().GetURL())
		if err != nil {
			return nil, err
		}

		clientConfig.BaseURL = url.Scheme + "://" + url.Host + "/v1"

		client = openai.NewClientWithConfig(clientConfig)
	}

	return &Engine{
		mode:         mode,
		config:       config,
		client:       client,
		execMessages: make([]openai.ChatCompletionMessage, 0),
		chatMessages: make([]openai.ChatCompletionMessage, 0),
		channel:      make(chan EngineChatStreamOutput),
		pipe:         "",
		running:      false,
	}, nil
}

func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mode = mode

	return e
}

func (e *Engine) GetMode() EngineMode {
	return e.mode
}

func (e *Engine) GetChannel() chan EngineChatStreamOutput {
	return e.channel
}

func (e *Engine) SetPipe(pipe string) *Engine {
	e.pipe = pipe

	return e
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineChatStreamOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}

	e.running = false

	return e
}

func (e *Engine) Clear() *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = []openai.ChatCompletionMessage{}
	} else {
		e.chatMessages = []openai.ChatCompletionMessage{}
	}

	return e
}

func (e *Engine) Reset() *Engine {
	e.execMessages = []openai.ChatCompletionMessage{}
	e.chatMessages = []openai.ChatCompletionMessage{}

	return e
}

func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     e.config.GetAIConfig().GetModel(),
			MaxTokens: e.config.GetAIConfig().GetMaxTokens(),
			Messages:  e.prepareCompletionMessages(),
		},
	)
	if err != nil {
		return nil, err
	}

	content := resp.Choices[0].Message.Content
	e.appendAssistantMessage(content)

	var output EngineExecOutput
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		re := regexp.MustCompile(`\{.*?\}`)
		match := re.FindString(content)
		if match != "" {
			err = json.Unmarshal([]byte(match), &output)
			if err != nil {
				return nil, err
			}
		} else {
			output = EngineExecOutput{
				Command:     "",
				Explanation: content,
				Executable:  false,
			}
		}
	}

	return &output, nil
}

func (e *Engine) ChatStreamCompletion(input string) error {
	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	req := openai.ChatCompletionRequest{
		Model:     e.config.GetAIConfig().GetModel(),
		MaxTokens: e.config.GetAIConfig().GetMaxTokens(),
		Messages:  e.prepareCompletionMessages(),
		Stream:    true,
	}

	stream, err := e.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	var output string

	for {
		if e.running {
			resp, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				executable := false
				if e.mode == ExecEngineMode {
					if !strings.HasPrefix(output, noexec) && !strings.Contains(output, "\n") {
						executable = true
					}
				}

				e.channel <- EngineChatStreamOutput{
					content:    "",
					last:       true,
					executable: executable,
				}
				e.running = false
				e.appendAssistantMessage(output)

				return nil
			}

			if err != nil {
				e.running = false
				return err
			}

			delta := resp.Choices[0].Delta.Content

			output += delta

			e.channel <- EngineChatStreamOutput{
				content: delta,
				last:    false,
			}

		} else {
			stream.Close()

			return nil
		}
	}
}

func (e *Engine) appendUserMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	}

	return e
}

func (e *Engine) appendAssistantMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	}

	return e
}

func (e *Engine) prepareCompletionMessages() []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: e.prepareSystemPrompt(),
		},
	}

	if e.pipe != "" {
		messages = append(
			messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: e.preparePipePrompt(),
			},
		)
	}

	if e.mode == ExecEngineMode {
		messages = append(messages, e.execMessages...)
	} else {
		messages = append(messages, e.chatMessages...)
	}

	return messages
}

func (e *Engine) preparePipePrompt() string {
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	return fmt.Sprintf("%s\n%s", bodyPart, e.prepareSystemPromptContextPart())
}

func (e *Engine) prepareSystemPromptExecPart() string {
	var sb strings.Builder

	sb.WriteString("You are Gogut, a powerful terminal assistant generating a JSON containing a command line for my input.\n")
	sb.WriteString("You will always reply using the following json structure: {\"cmd\":\"the command\", \"exp\": \"some explanation\", \"exec\": true}.\n")
	sb.WriteString("Your answer will always only contain the json structure, never add any advice or supplementary detail or information, even if I asked the same question before.\n")
	sb.WriteString("The field cmd will contain a single line command (don't use new lines, use separators like && and ; instead).\n")
	sb.WriteString("The field exp will contain a short explanation of the command if you managed to generate an executable command, otherwise it will contain the reason of your failure.\n")
	sb.WriteString("The field exec will contain true if you managed to generate an executable command, false otherwise.\n")
	sb.WriteString("\n")
	sb.WriteString("Examples:\n")
	sb.WriteString("Me: list all files in my home dir\n")
	sb.WriteString("Gogut: {\"cmd\":\"ls ~\", \"exp\": \"list all files in your home dir\", \"exec\": true}\n")
	sb.WriteString("Me: list all pods of all namespaces\n")
	sb.WriteString("Gogut: {\"cmd\":\"kubectl get pods --all-namespaces\", \"exp\": \"list pods from all k8s namespaces\", \"exec\": true}\n")
	sb.WriteString("Me: how are you ?\n")
	sb.WriteString("Gogut: {\"cmd\":\"\", \"exp\": \"I'm good thanks but I cannot generate a command for this. Use the chat mode to discuss.\", \"exec\": false}")

	return sb.String()
}

func (e *Engine) prepareSystemPromptChatPart() string {
	var sb strings.Builder

	sb.WriteString("You are Gogut, a powerful terminal assistant created by github.com/bmichalkiewicz.\n")
	sb.WriteString("You will answer in the most helpful possible way.\n")
	sb.WriteString("Always format your answer in markdown format.\n\n")
	sb.WriteString("For example:\n")
	sb.WriteString("Me: What is 2+2 ?\n")
	sb.WriteString("Gogut: The answer for `2+2` is `4`\n")
	sb.WriteString("Me: +2 again ?\n")
	sb.WriteString("Gogut: The answer is `6`\n")

	return sb.String()
}

func (e *Engine) prepareSystemPromptContextPart() string {
	var sb strings.Builder

	sb.WriteString("My context: ")

	if e.config.GetSystemConfig().GetOperatingSystem() != facts.UnknownOperatingSystem {
		sb.WriteString(fmt.Sprintf("my operating system is %s, ", e.config.GetSystemConfig().GetOperatingSystem().String()))
	}
	if e.config.GetSystemConfig().GetDistribution() != "" {
		sb.WriteString(fmt.Sprintf("my distribution is %s, ", e.config.GetSystemConfig().GetDistribution()))
	}
	if e.config.GetSystemConfig().GetHomeDirectory() != "" {
		sb.WriteString(fmt.Sprintf("my home directory is %s, ", e.config.GetSystemConfig().GetHomeDirectory()))
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		sb.WriteString(fmt.Sprintf("my shell is %s, ", e.config.GetSystemConfig().GetShell()))
	}
	if e.config.GetSystemConfig().GetEditor() != "" {
		sb.WriteString(fmt.Sprintf("my editor is %s, ", e.config.GetSystemConfig().GetEditor()))
	}
	sb.WriteString("take this into account. ")

	if e.config.GetUserConfig().GetPreferences() != "" {
		sb.WriteString(fmt.Sprintf("Also, %s.", e.config.GetUserConfig().GetPreferences()))
	}

	return sb.String()
}
