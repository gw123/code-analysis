package cmd

import (
	"codetest"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dir       string
	apiToken  string
	outputDir string
)

// analyzeCmd 定义了分析命令
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze code in the specified directory using AI",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(dir, apiToken)
	},
}

// init 函数用于设置分析命令的参数
func init() {
	rootCmd.AddCommand(analyzeCmd) // 将子命令添加到根命令

	// 通过 flag 接受目录和 API token
	analyzeCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to analyze (required)")
	analyzeCmd.Flags().StringVarP(&apiToken, "token", "t", "", "API token for AI analysis (required)")
	analyzeCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "./result", "总结文件输出地方")

	// 必须参数检查
	err := analyzeCmd.MarkFlagRequired("dir")
	if err != nil {
		log.Println("Error: dir flag is required", err)
		return
	}
	err = analyzeCmd.MarkFlagRequired("token")
	if err != nil {
		log.Println("Error: token flag is required", err)
		return
	}
}

// run 主要逻辑
func run(directory, token string) error {
	aiClient := code.NewChatGPTClient(token)
	var count int

	// 遍历目录并处理每个文件
	err := code.WalkDir(directory, func(path string) {
		processFile(path, aiClient)
		count++
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	fmt.Printf("Processed %d files\n", count)
	return nil
}

// 处理单个文件
func processFile(path string, aiClient *code.ChatGPTClient) {
	fmt.Println("Processing file:", path)

	// 读取文件内容
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", path, err)
		return
	}

	// 调用 AI 进行代码分析
	rawAiResponse, yamlResult, err := aiClient.AIAnalysisCode(path, string(fileContent))
	if err != nil {
		log.Printf("AI analysis failed for %s: %v\n", path, err)
		return
	}

	// 生成文件名并保存分析结果
	if err := saveAIResult(path, rawAiResponse); err != nil {
		log.Printf("Failed to save AI result for %s: %v\n", path, err)
		return
	}

	// 更新总结文件
	if err := updateSummaryFile(path, &yamlResult); err != nil {
		log.Printf("Failed to update summary for %s: %v\n", path, err)
		return
	}
}

// 保存 AI 分析结果到文件
func saveAIResult(path, rawAiResponse string) error {
	resultPath := filepath.Join(outputDir, strings.ReplaceAll(path, "/", "|")+".yaml")
	if err := os.WriteFile(resultPath, []byte(rawAiResponse), 0644); err != nil {
		return fmt.Errorf("error writing result file: %v", err)
	}
	return nil
}

// 更新总结文件
func updateSummaryFile(path string, yamlResult *code.ParsedYAML) error {
	var strBuilder strings.Builder

	strBuilder.WriteString(fmt.Sprintf("文件名: %s\n", path))
	strBuilder.WriteString(fmt.Sprintf("功能: %s\n", yamlResult.FunctionDescription))
	strBuilder.WriteString(fmt.Sprintf("包名: %s\n", yamlResult.FileInfo.PackageName))
	strBuilder.WriteString("依赖导入项目: ")
	strBuilder.WriteString(strings.Join(yamlResult.FileInfo.Imports, ","))
	strBuilder.WriteString("\n---\n")

	// 追加写入总结文件
	file, err := os.OpenFile(outputDir+"/all.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open summary file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(strBuilder.String()); err != nil {
		return fmt.Errorf("failed to write to summary file: %v", err)
	}
	return nil
}
