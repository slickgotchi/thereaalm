package ai

// llm API interaction

import "log"

// SendToLLM sends Gotchi data to the LLM and gets a decision.
func SendToLLM(context string) string {
    log.Printf("Sending context to LLM: %s\n", context)
    return "LLM Response: Investigate nearby object."
}