package anthropic

import "github.com/maximhq/bifrost/core/schemas"

// ToBifrostCountTokensResponse converts the Anthropic count tokens response into Bifrost format.
func (resp *AnthropicCountTokensResponse) ToBifrostCountTokensResponse(provider schemas.ModelProvider) *schemas.BifrostCountTokensResponse {
	if resp == nil {
		return nil
	}

	var inputDetails *schemas.ResponsesResponseInputTokens
	if resp.InputTokensDetails != nil {
		inputDetails = &schemas.ResponsesResponseInputTokens{
			AudioTokens:  resp.InputTokensDetails.AudioTokens,
			CachedTokens: resp.InputTokensDetails.CachedTokens,
		}
	}

	return &schemas.BifrostCountTokensResponse{
		Object:             resp.Object,
		Model:              resp.Model,
		InputTokens:        resp.InputTokens,
		InputTokensDetails: inputDetails,
		OutputTokens:       resp.OutputTokens,
		TotalTokens:        resp.TotalTokens,
		Usage:              convertAnthropicUsageToResponsesUsage(resp.Usage),
		ExtraFields: schemas.BifrostResponseExtraFields{
			Provider:       provider,
			ModelRequested: resp.Model,
		},
	}
}

// ToAnthropicCountTokensRequest converts a Bifrost count tokens request into Anthropic's message format.
func ToAnthropicCountTokensRequest(bifrostReq *schemas.BifrostCountTokensRequest) *AnthropicMessageRequest {
	if bifrostReq == nil {
		return nil
	}

	responsesReq := &schemas.BifrostResponsesRequest{
		Provider:  bifrostReq.Provider,
		Model:     bifrostReq.Model,
		Input:     bifrostReq.Input,
		Params:    bifrostReq.Params,
		Fallbacks: bifrostReq.Fallbacks,
	}

	return ToAnthropicResponsesRequest(responsesReq)
}

func convertAnthropicUsageToResponsesUsage(usage *AnthropicUsage) *schemas.ResponsesResponseUsage {
	if usage == nil {
		return nil
	}

	respUsage := &schemas.ResponsesResponseUsage{
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		TotalTokens:  usage.InputTokens + usage.OutputTokens,
	}

	if usage.CacheReadInputTokens > 0 {
		respUsage.InputTokensDetails = &schemas.ResponsesResponseInputTokens{
			CachedTokens: usage.CacheReadInputTokens,
		}
	}
	if usage.CacheCreationInputTokens > 0 {
		respUsage.OutputTokensDetails = &schemas.ResponsesResponseOutputTokens{
			CachedTokens: usage.CacheCreationInputTokens,
		}
	}

	return respUsage
}
