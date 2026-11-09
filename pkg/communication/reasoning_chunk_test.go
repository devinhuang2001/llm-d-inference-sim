package communication

import (
    "strings"
    "testing"
)

func TestConcatenatedSSEChunkKeepsFramesSeparate(t *testing.T) {
    frames := concatenatedSSEChunk{
        &jsonDataChunk{data: map[string]any{"delta": map[string]any{"reasoning": "thinking"}}},
        &jsonDataChunk{data: map[string]any{"delta": map[string]any{"content": "answer"}}},
    }
    payload, err := frames.SSEBytes()
    if err != nil { t.Fatal(err) }
    parts := strings.Split(strings.TrimSpace(string(payload)), "\n\n")
    if len(parts) != 2 { t.Fatalf("want two frames, got %q", payload) }
    if !strings.Contains(parts[0], "\"reasoning\"") || strings.Contains(parts[0], "\"content\"") { t.Fatalf("wrong reasoning frame: %q", parts[0]) }
    if !strings.Contains(parts[1], "\"content\"") || strings.Contains(parts[1], "\"reasoning\"") { t.Fatalf("wrong content frame: %q", parts[1]) }
}
