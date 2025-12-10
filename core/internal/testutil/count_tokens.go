package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"
)

// RunCountTokenTest validates the CountTokens API for the configured provider/model.
// It sends a simple prompt as Responses messages and asserts token counts and metadata.
func RunCountTokenTest(t *testing.T, client *bifrost.Bifrost, ctx context.Context, testConfig ComprehensiveTestConfig) {
	if !testConfig.Scenarios.CountTokens {
		t.Logf("Count token not supported for provider %s", testConfig.Provider)
		return
	}

	t.Run("CountTokens", func(t *testing.T) {
		if os.Getenv("SKIP_PARALLEL_TESTS") != "true" {
			t.Parallel()
		}

		messages := []schemas.ResponsesMessage{
			CreateBasicResponsesMessage("Hello! What's the capital of France?"),
		}

		retryConfig := GetTestRetryConfigForScenario("CountTokens", testConfig)

		countTokensReq := &schemas.BifrostCountTokensRequest{
			Provider:  testConfig.Provider,
			Model:     testConfig.ChatModel,
			Input:     messages,
			Params:    &schemas.ResponsesParameters{Temperature: bifrost.Ptr(0.2)},
			Fallbacks: testConfig.Fallbacks,
		}

		retryContext := TestRetryContext{
			ScenarioName: "CountTokens",
			ExpectedBehavior: map[string]interface{}{
				"should_return_token_counts": true,
			},
			TestMetadata: map[string]interface{}{
				"provider": testConfig.Provider,
				"model":    testConfig.ChatModel,
			},
		}

		// Create CountTokens retry config
		countTokensRetryConfig := CountTokensRetryConfig{
			MaxAttempts: retryConfig.MaxAttempts,
			BaseDelay:   retryConfig.BaseDelay,
			MaxDelay:    retryConfig.MaxDelay,
			Conditions:  []CountTokensRetryCondition{},
			OnRetry:     retryConfig.OnRetry,
			OnFinalFail: retryConfig.OnFinalFail,
		}

		// Validation function
		validateCountTokens := func(resp *schemas.BifrostCountTokensResponse) error {
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			if resp.Model != countTokensReq.Model {
				return fmt.Errorf("model mismatch: got %s want %s", resp.Model, countTokensReq.Model)
			}
			if resp.InputTokens <= 0 {
				return fmt.Errorf("input_tokens should be > 0, got %d", resp.InputTokens)
			}
			if resp.TotalTokens < resp.InputTokens {
				return fmt.Errorf("total_tokens (%d) should be >= input_tokens (%d)", resp.TotalTokens, resp.InputTokens)
			}
			if resp.Usage == nil {
				return fmt.Errorf("usage should be populated")
			}
			if resp.Usage.TotalTokens != resp.TotalTokens {
				return fmt.Errorf("usage.total_tokens mismatch: got %d want %d", resp.Usage.TotalTokens, resp.TotalTokens)
			}
			if resp.ExtraFields.RequestType != schemas.CountTokensRequest {
				return fmt.Errorf("request type not set: got %s", resp.ExtraFields.RequestType)
			}
			if resp.ExtraFields.Provider != testConfig.Provider {
				return fmt.Errorf("provider not set on extra fields: got %s want %s", resp.ExtraFields.Provider, testConfig.Provider)
			}
			return nil
		}

		// Use retry framework
		countTokensResp, countTokensErr := WithCountTokensTestRetry(
			func() (*schemas.BifrostCountTokensResponse, *schemas.BifrostError) {
				return client.CountTokensRequest(ctx, countTokensReq)
			},
			validateCountTokens,
			countTokensRetryConfig,
			retryContext,
			t,
		)

		if countTokensErr != nil {
			t.Fatalf("❌ CountTokens request failed: %s", GetErrorMessage(countTokensErr))
		}
		if countTokensResp == nil {
			t.Fatal("❌ CountTokens response is nil")
		}

		// All validations are handled in the validation function
		t.Logf("✅ CountTokens test passed: input=%d, total=%d", countTokensResp.InputTokens, countTokensResp.TotalTokens)
	})
}
