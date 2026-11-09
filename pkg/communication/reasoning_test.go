package communication

import (
	"strings"
	"testing"

	"github.com/llm-d/llm-d-inference-sim/pkg/api"
	"github.com/llm-d/llm-d-inference-sim/pkg/endpoint"
)

type reasoningTestContext struct {
	endpoint.ResponseContext
	percentage int
	tools []api.ToolCall
	sendImage bool
	audio bool
}

func (c reasoningTestContext) ReasoningPercentage() int { return c.percentage }
func (c reasoningTestContext) ToolCalls() []api.ToolCall { return c.tools }
func (c reasoningTestContext) SendImage() bool { return c.sendImage }
func (c reasoningTestContext) HasAudioOutput() bool { return c.audio }

func TestChatReasoningBudget(t *testing.T) {
	c := reasoningTestContext{percentage: 25}
	for _, tt := range []struct { tokens, expected int }{
		{0, 0}, {3, 0}, {4, 1}, {10, 2}, {100, 25},
	} {
		if got := chatReasoningTokens(c, tt.tokens); got != tt.expected {
			t.Errorf("count=%d: got %d, want %d", tt.tokens, got, tt.expected)
		}
	}
	c.tools = []api.ToolCall{{}}
	if got := chatReasoningTokens(c, 100); got != 0 {
		t.Errorf("tool calls: got %d, want 0", got)
	}
	c.tools = nil
	c.sendImage = true
	if got := chatReasoningTokens(c, 100); got != 0 {
		t.Errorf("image: got %d, want 0", got)
	}
}

func TestMixedReasoningFramesRemainSeparate(t *testing.T) {
	chunks := concatenatedSSEChunk{
		&jsonDataChunk{data: map[string]any{"delta": map[string]any{"reasoning": "think"}}},
		&jsonDataChunk{data: map[string]any{"delta": map[string]any{"content": "answer"}}},
	}
	wire, err := chunks.SSEBytes()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(wire), "reasoning") || !strings.Contains(string(wire), "content") {
		t.Fatalf("missing split fields: %s", wire)
	}
	if strings.Count(string(wire), "data: ") != 2 {
		t.Fatalf("expected two SSE frames: %s", wire)
	}
}
