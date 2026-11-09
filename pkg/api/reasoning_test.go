package api

import (
    "encoding/json"
    "testing"
)

func TestReasoningMessageWireFormat(t *testing.T) {
    withReasoning, err := json.Marshal(Message{Role: RoleAssistant, Reasoning: "think"})
    if err != nil { t.Fatal(err) }
    var decoded map[string]any
    if err := json.Unmarshal(withReasoning, &decoded); err != nil { t.Fatal(err) }
    if decoded["reasoning"] != "think" { t.Fatalf("missing reasoning field: %s", withReasoning) }
    if value, ok := decoded["content"]; !ok || value != nil { t.Fatalf("reasoning-only message should have null content: %s", withReasoning) }

    withoutReasoning, err := json.Marshal(Message{Role: RoleAssistant, Content: ChatComplContent{Raw: "answer"}})
    if err != nil { t.Fatal(err) }
    decoded = map[string]any{}
    if err := json.Unmarshal(withoutReasoning, &decoded); err != nil { t.Fatal(err) }
    if _, ok := decoded["reasoning"]; ok { t.Fatalf("default output contains reasoning: %s", withoutReasoning) }
}
