package schemas

// BifrostCountTokensRequest represents a request to count tokens for a given model/input pair.
type BifrostCountTokensRequest struct {
	Provider       ModelProvider        `json:"provider"`
	Model          string               `json:"model"`
	Input          []ResponsesMessage   `json:"input,omitempty"`
	Params         *ResponsesParameters `json:"params,omitempty"`
	Fallbacks      []Fallback           `json:"fallbacks,omitempty"`
	RawRequestBody []byte               `json:"-"` // set bifrost-use-raw-request-body to true in ctx to use the raw request body. Bifrost will directly send this to the downstream provider.
}

func (r *BifrostCountTokensRequest) GetRawRequestBody() []byte {
	return r.RawRequestBody
}

// BifrostCountTokensResponse captures token counts for a provided input.
type BifrostCountTokensResponse struct {
	Object             string                        `json:"object,omitempty"`
	Model              string                        `json:"model"`
	InputTokens        int                           `json:"input_tokens"`
	InputTokensDetails *ResponsesResponseInputTokens `json:"input_tokens_details,omitempty"`
	OutputTokens       int                           `json:"output_tokens,omitempty"`
	TotalTokens        int                           `json:"total_tokens"`
	Usage              *ResponsesResponseUsage       `json:"usage,omitempty"`
	ExtraFields        BifrostResponseExtraFields    `json:"extra_fields"`
}
