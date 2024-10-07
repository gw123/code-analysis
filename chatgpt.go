package code

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

// Step1FileInfo 结构体表示文件信息
type Step1FileInfo struct {
	File        string
	Why         string
	ParseResult string
}

// ChatGPTClient 结构体封装 ChatGPT 客户端
type ChatGPTClient struct {
	client *openai.Client
}

// NewChatGPTClient 创建新的 ChatGPTClient
func NewChatGPTClient(apiKey string) *ChatGPTClient {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://api.chatanywhere.tech/v1"
	return &ChatGPTClient{
		client: openai.NewClientWithConfig(cfg),
	}
}

// getChatGPTResponse 调用 ChatGPT API 并返回回复
func (c *ChatGPTClient) getChatGPTResponse(prompt string) (string, error) {
	ctx := context.Background()
	// 构造请求消息
	req := openai.ChatCompletionRequest{
		Temperature: 0,
		Model:       openai.GPT4oMini, // 使 用 GPT-3.5 Turbo 模型
		//Model: openai.CodexCodeDavinci002, // 使用 GPT-3.5 Turbo 模型
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	// 调用 OpenAI API 获取回复
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ChatGPT request failed: %v", err)
	}

	// 返回模型的回复内容
	return resp.Choices[0].Message.Content, nil
}

func (c *ChatGPTClient) AIAnalysisCode(filename, code string) (string, ParsedYAML, error) {
	response, err := c.getChatGPTResponse(buildFileAnalysisPrompt(filename, code))
	if err != nil {
		return "", ParsedYAML{}, err
	}
	response = strings.TrimSpace(response)
	response = strings.TrimLeft(response, "```yaml")
	response = strings.TrimLeft(response, "\n")
	response = strings.TrimRight(response, "```")

	// 修正大模型的错误
	{
		response = strings.ReplaceAll(response, "structs: []\n", "")
		response = strings.ReplaceAll(response, "structs: ''\n", "")
		response = strings.ReplaceAll(response, "constants: ''\n", "")
		response = strings.ReplaceAll(response, "constants: []\n", "")
		response = strings.ReplaceAll(response, "constants: []\n", "")
		response = strings.ReplaceAll(response, "interfaces: ''\n", "")
		response = strings.ReplaceAll(response, "interfaces: []\n", "")
		response = strings.ReplaceAll(response, "params: ''\n", "")
		response = strings.ReplaceAll(response, "return_values: ''\n", "")
		response = strings.ReplaceAll(response, "- []\n", "")

		regex := regexp.MustCompile(`- (\w+): \*(.*)`)
		response = regex.ReplaceAllString(response, `$1: '*$2'`)
	}

	var parsedData ParsedYAML
	err = yaml.Unmarshal([]byte(response), &parsedData)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return response, parsedData, nil
	}

	return response, parsedData, nil
}

func (c *ChatGPTClient) AIQuestion(summaryContent, question, helpInfo string) ([]string, error) {

	step1Response, err := c.getChatGPTResponse(buildQuestionRelFilesPrompt(question, summaryContent))
	if err != nil {
		return nil, err
	}

	step1Response = strings.TrimSpace(step1Response)
	step1Response = strings.TrimLeft(step1Response, "```yaml")
	step1Response = strings.TrimRight(step1Response, "```")

	var step1FileInfos []*Step1FileInfo
	err = yaml.Unmarshal([]byte(step1Response), &step1FileInfos)
	if err != nil {
		fmt.Println("Step1FileInfo Error parsing YAML2:", err)
		fmt.Println(step1Response)
		return nil, err
	}

	fmt.Println("----------需要召回的文件列表-------------")
	for _, step1FileInfo := range step1FileInfos {
		fmt.Println(step1FileInfo.File)
		fmt.Println(step1FileInfo.Why)
	}

	for _, step1FileInfo := range step1FileInfos {
		fileContent, err := os.ReadFile(step1FileInfo.File)
		if err != nil {
			return nil, err
		}

		response, err := c.getChatGPTResponse(buildQuestionRelFilesParsePrompt(question, step1Response, step1FileInfo.File, string(fileContent)))
		if err != nil {
			return nil, err
		}
		fmt.Println("-----", step1FileInfo.File, "分析结果")
		fmt.Println(response)
		step1FileInfo.ParseResult = response
	}

	answerPromptBuilder := buildFinalAnswerPrompt(question, helpInfo)
	for _, step1FileInfo := range step1FileInfos {
		answerPromptBuilder.WriteString(string(step1FileInfo.ParseResult))
	}

	response, err := c.getChatGPTResponse(answerPromptBuilder.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(response)
	return nil, nil
}

func (c *ChatGPTClient) GenNodeDoc(nodeName, fileContent string) (string, error) {
	prompt := strings.Builder{}
	prompt.WriteString("生成节点使用文档\n节点名称：auth")
	prompt.WriteString(nodeName)
	prompt.WriteString("\n")
	prompt.WriteString(GenNodeHelpInfo())
	prompt.WriteString(GenCodeUseDocHelpInfo())

	prompt.WriteString(fileContent)
	fmt.Println("######################")
	fmt.Println(prompt.String())
	fmt.Println("######################")
	return c.getChatGPTResponse(prompt.String())
}

func (c *ChatGPTClient) GenWorkflowYaml(workflowUsage, allNodeUsage string) (string, error) {
	return c.getChatGPTResponse(GenWorkflowYaml(workflowUsage, allNodeUsage))
}
