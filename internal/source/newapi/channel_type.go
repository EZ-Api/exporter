// Package newapi provides channel type enum mapping.
// Reference: SPEC_newapi_migration_tool.md Appendix B
package newapi

// ChannelType represents the channel type enum in New API.
type ChannelType int

// Channel type constants from New API constant/channel.go
const (
	ChannelTypeUnknown        ChannelType = 0
	ChannelTypeOpenAI         ChannelType = 1
	ChannelTypeMidjourney     ChannelType = 2
	ChannelTypeAzure          ChannelType = 3
	ChannelTypeOllama         ChannelType = 4
	ChannelTypeMidjourneyPlus ChannelType = 5
	ChannelTypeOpenAIMax      ChannelType = 6
	ChannelTypeOhMyGPT        ChannelType = 7
	ChannelTypeCustom         ChannelType = 8
	ChannelTypeAILS           ChannelType = 9
	ChannelTypeAIProxy        ChannelType = 10
	ChannelTypePaLM           ChannelType = 11
	ChannelTypeAPI2GPT        ChannelType = 12
	ChannelTypeAIGC2D         ChannelType = 13
	ChannelTypeAnthropic      ChannelType = 14
	ChannelTypeBaidu          ChannelType = 15
	ChannelTypeZhipu          ChannelType = 16
	ChannelTypeAli            ChannelType = 17
	ChannelTypeXunfei         ChannelType = 18
	ChannelType360            ChannelType = 19
	ChannelTypeOpenRouter     ChannelType = 20
	ChannelTypeAIProxyLibrary ChannelType = 21
	ChannelTypeFastGPT        ChannelType = 22
	ChannelTypeTencent        ChannelType = 23
	ChannelTypeGemini         ChannelType = 24
	ChannelTypeMoonshot       ChannelType = 25
	ChannelTypeZhipuV4        ChannelType = 26
	ChannelTypePerplexity     ChannelType = 27
	ChannelTypeLingYiWanWu    ChannelType = 31
	ChannelTypeAws            ChannelType = 33
	ChannelTypeCohere         ChannelType = 34
	ChannelTypeMiniMax        ChannelType = 35
	ChannelTypeSunoAPI        ChannelType = 36
	ChannelTypeDify           ChannelType = 37
	ChannelTypeJina           ChannelType = 38
	ChannelTypeCloudflare     ChannelType = 39
	ChannelTypeSiliconFlow    ChannelType = 40
	ChannelTypeVertexAI       ChannelType = 41
	ChannelTypeMistral        ChannelType = 42
	ChannelTypeDeepSeek       ChannelType = 43
	ChannelTypeMokaAI         ChannelType = 44
	ChannelTypeVolcEngine     ChannelType = 45
	ChannelTypeBaiduV2        ChannelType = 46
	ChannelTypeXinference     ChannelType = 47
	ChannelTypeXai            ChannelType = 48
	ChannelTypeCoze           ChannelType = 49
	ChannelTypeKling          ChannelType = 50
	ChannelTypeJimeng         ChannelType = 51
	ChannelTypeVidu           ChannelType = 52
	ChannelTypeSubmodel       ChannelType = 53
	ChannelTypeDoubaoVideo    ChannelType = 54
	ChannelTypeSora           ChannelType = 55
	ChannelTypeReplicate      ChannelType = 56
)

// channelTypeMapping maps New API channel type int to EZ-API provider type string.
var channelTypeMapping = map[ChannelType]string{
	ChannelTypeUnknown:        "custom",
	ChannelTypeOpenAI:         "openai",
	ChannelTypeMidjourney:     "midjourney",
	ChannelTypeAzure:          "azure",
	ChannelTypeOllama:         "ollama",
	ChannelTypeMidjourneyPlus: "midjourney",
	ChannelTypeOpenAIMax:      "openai",
	ChannelTypeOhMyGPT:        "openai",
	ChannelTypeCustom:         "custom",
	ChannelTypeAILS:           "openai",
	ChannelTypeAIProxy:        "openai",
	ChannelTypePaLM:           "palm",
	ChannelTypeAPI2GPT:        "openai",
	ChannelTypeAIGC2D:         "openai",
	ChannelTypeAnthropic:      "anthropic",
	ChannelTypeBaidu:          "baidu",
	ChannelTypeZhipu:          "zhipu",
	ChannelTypeAli:            "ali",
	ChannelTypeXunfei:         "xunfei",
	ChannelType360:            "360",
	ChannelTypeOpenRouter:     "openrouter",
	ChannelTypeAIProxyLibrary: "openai",
	ChannelTypeFastGPT:        "openai",
	ChannelTypeTencent:        "tencent",
	ChannelTypeGemini:         "gemini",
	ChannelTypeMoonshot:       "moonshot",
	ChannelTypeZhipuV4:        "zhipu",
	ChannelTypePerplexity:     "perplexity",
	ChannelTypeLingYiWanWu:    "lingyiwanwu",
	ChannelTypeAws:            "aws",
	ChannelTypeCohere:         "cohere",
	ChannelTypeMiniMax:        "minimax",
	ChannelTypeSunoAPI:        "suno",
	ChannelTypeDify:           "dify",
	ChannelTypeJina:           "jina",
	ChannelTypeCloudflare:     "cloudflare",
	ChannelTypeSiliconFlow:    "siliconflow",
	ChannelTypeVertexAI:       "vertex",
	ChannelTypeMistral:        "mistral",
	ChannelTypeDeepSeek:       "deepseek",
	ChannelTypeMokaAI:         "mokaai",
	ChannelTypeVolcEngine:     "volcengine",
	ChannelTypeBaiduV2:        "baidu",
	ChannelTypeXinference:     "xinference",
	ChannelTypeXai:            "xai",
	ChannelTypeCoze:           "coze",
	ChannelTypeKling:          "kling",
	ChannelTypeJimeng:         "jimeng",
	ChannelTypeVidu:           "vidu",
	ChannelTypeSubmodel:       "submodel",
	ChannelTypeDoubaoVideo:    "doubao",
	ChannelTypeSora:           "openai",
	ChannelTypeReplicate:      "replicate",
}

// channelTypeDisplayNames provides human-readable names for channel types.
var channelTypeDisplayNames = map[ChannelType]string{
	ChannelTypeUnknown:        "Unknown",
	ChannelTypeOpenAI:         "OpenAI",
	ChannelTypeMidjourney:     "Midjourney",
	ChannelTypeAzure:          "Azure",
	ChannelTypeOllama:         "Ollama",
	ChannelTypeMidjourneyPlus: "MidjourneyPlus",
	ChannelTypeOpenAIMax:      "OpenAIMax",
	ChannelTypeOhMyGPT:        "OhMyGPT",
	ChannelTypeCustom:         "Custom",
	ChannelTypeAILS:           "AILS",
	ChannelTypeAIProxy:        "AIProxy",
	ChannelTypePaLM:           "PaLM",
	ChannelTypeAPI2GPT:        "API2GPT",
	ChannelTypeAIGC2D:         "AIGC2D",
	ChannelTypeAnthropic:      "Anthropic",
	ChannelTypeBaidu:          "Baidu",
	ChannelTypeZhipu:          "Zhipu",
	ChannelTypeAli:            "Ali",
	ChannelTypeXunfei:         "Xunfei",
	ChannelType360:            "360",
	ChannelTypeOpenRouter:     "OpenRouter",
	ChannelTypeAIProxyLibrary: "AIProxyLibrary",
	ChannelTypeFastGPT:        "FastGPT",
	ChannelTypeTencent:        "Tencent",
	ChannelTypeGemini:         "Gemini",
	ChannelTypeMoonshot:       "Moonshot",
	ChannelTypeZhipuV4:        "ZhipuV4",
	ChannelTypePerplexity:     "Perplexity",
	ChannelTypeLingYiWanWu:    "LingYiWanWu",
	ChannelTypeAws:            "AWS",
	ChannelTypeCohere:         "Cohere",
	ChannelTypeMiniMax:        "MiniMax",
	ChannelTypeSunoAPI:        "SunoAPI",
	ChannelTypeDify:           "Dify",
	ChannelTypeJina:           "Jina",
	ChannelTypeCloudflare:     "Cloudflare",
	ChannelTypeSiliconFlow:    "SiliconFlow",
	ChannelTypeVertexAI:       "VertexAI",
	ChannelTypeMistral:        "Mistral",
	ChannelTypeDeepSeek:       "DeepSeek",
	ChannelTypeMokaAI:         "MokaAI",
	ChannelTypeVolcEngine:     "VolcEngine",
	ChannelTypeBaiduV2:        "BaiduV2",
	ChannelTypeXinference:     "Xinference",
	ChannelTypeXai:            "xAI",
	ChannelTypeCoze:           "Coze",
	ChannelTypeKling:          "Kling",
	ChannelTypeJimeng:         "Jimeng",
	ChannelTypeVidu:           "Vidu",
	ChannelTypeSubmodel:       "Submodel",
	ChannelTypeDoubaoVideo:    "DoubaoVideo",
	ChannelTypeSora:           "Sora",
	ChannelTypeReplicate:      "Replicate",
}

// ToProviderType converts a New API channel type to EZ-API provider type string.
// Returns ("custom", false) if the type is unknown.
func (t ChannelType) ToProviderType() (string, bool) {
	providerType, ok := channelTypeMapping[t]
	if !ok {
		return "custom", false
	}
	return providerType, true
}

// DisplayName returns the human-readable name for the channel type.
func (t ChannelType) DisplayName() string {
	name, ok := channelTypeDisplayNames[t]
	if !ok {
		return "Unknown"
	}
	return name
}

// MapChannelType maps an integer channel type to EZ-API provider type.
// Returns the mapped type and a boolean indicating if mapping was successful.
func MapChannelType(typeID int) (string, bool) {
	return ChannelType(typeID).ToProviderType()
}
