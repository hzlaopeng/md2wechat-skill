package image

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

// DashScopeProvider 阿里云 DashScope 图片生成服务提供商
type DashScopeProvider struct {
	apiKey  string
	baseURL string
	model   string
	size    string
	client  *http.Client
}

// NewDashScopeProvider 创建 DashScope Provider
func NewDashScopeProvider(cfg *config.Config) (*DashScopeProvider, error) {
	model := cfg.ImageModel
	if model == "" {
		model = "z-image-turbo" // 默认模型
	}

	size := cfg.ImageSize
	if size == "" {
		size = "1024*1024" // 阿里云格式
	}

	baseURL := cfg.ImageAPIBase
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/api/v1"
	}

	return &DashScopeProvider{
		apiKey:  cfg.ImageAPIKey,
		baseURL: baseURL,
		model:   model,
		size:    size,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Name 返回提供商名称
func (p *DashScopeProvider) Name() string {
	return "DashScope"
}

// Generate 生成图片
func (p *DashScopeProvider) Generate(ctx context.Context, prompt string) (*GenerateResult, error) {
	// 构建请求体 - 阿里云 DashScope 格式 (使用 messages)
	reqBody := map[string]any{
		"model": p.model,
		"input": map[string]any{
			"messages": []map[string]any{
				{
					"role": "user",
					"content": []map[string]string{
						{"text": prompt},
					},
				},
			},
		},
		"parameters": map[string]any{
			"size": p.size,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "marshal_error",
			Message:  "构建请求失败",
			Original: err,
		}
	}

	// 发送请求
	url := p.baseURL + "/services/aigc/multimodal-generation/generation"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "request_error",
			Message:  "创建请求失败",
			Original: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "network_error",
			Message:  "网络请求失败，请检查网络连接",
			Hint:     "确保网络正常且 API 地址正确",
			Original: err,
		}
	}
	defer resp.Body.Close()

	// 处理错误响应
	if resp.StatusCode != http.StatusOK {
		return nil, p.handleErrorResponse(resp)
	}

	// 解析响应 - 阿里云 DashScope 格式
	var result struct {
		Output struct {
			Choices []struct {
				Message struct {
					Content []struct {
						Image string `json:"image"`
						Text  string `json:"text"`
					} `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"output"`
		Usage struct {
			ImageCount int `json:"image_count"`
		} `json:"usage"`
		RequestID string `json:"request_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "decode_error",
			Message:  "解析响应失败",
			Original: err,
		}
	}

	if len(result.Output.Choices) == 0 || len(result.Output.Choices[0].Message.Content) == 0 {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "no_image",
			Message:  "未生成图片",
			Hint:     "可能提示词不符合规范，请尝试修改提示词",
		}
	}

	imageURL := result.Output.Choices[0].Message.Content[0].Image
	if imageURL == "" {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "no_image",
			Message:  "未生成图片",
			Hint:     "API 返回的图片 URL 为空",
		}
	}

	return &GenerateResult{
		URL:   imageURL,
		Model: p.model,
		Size:  p.size,
	}, nil
}

// handleErrorResponse 处理错误响应
func (p *DashScopeProvider) handleErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		RequestID string `json:"request_id"`
	}

	_ = json.Unmarshal(body, &errResp)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unauthorized",
			Message:  "API Key 无效或已过期",
			Hint:     "请检查配置文件中的 api.image_key 是否正确",
			Original: fmt.Errorf("status 401: %s", string(body)),
		}
	case http.StatusTooManyRequests:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "rate_limit",
			Message:  "请求频率超限，请稍后再试",
			Hint:     "DashScope API 有请求频率限制，请等待一段时间后重试",
			Original: fmt.Errorf("status 429: %s", string(body)),
		}
	case http.StatusBadRequest:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "bad_request",
			Message:  fmt.Sprintf("请求参数错误: %s", errResp.Message),
			Hint:     "请检查图片尺寸、模型名称等参数是否正确",
			Original: fmt.Errorf("status 400: %s", string(body)),
		}
	default:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unknown",
			Message:  fmt.Sprintf("API 返回错误 (HTTP %d): %s", resp.StatusCode, errResp.Message),
			Hint:     "请稍后重试，或查看 DashScope 控制台状态",
			Original: fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		}
	}
}
