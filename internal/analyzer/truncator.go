package analyzer

import (
	"context"

	"github.com/priyansh-dimri/argus/pkg/logger"
)

func TruncateLog(ctx context.Context, client AIClient, log string, maxTokens int) string {
	if len(log) < maxTokens {
		return log
	}

	count, err := client.CountTokens(ctx, log)
	if err != nil {
		logger.Warn("Failed to count tokens during truncation", "error", err)
		safeCharLimit := maxTokens * 3
		if len(log) > safeCharLimit {
			return log[:safeCharLimit] + "...[TRUNCATED_SAFE_MODE]"
		}
		return log
	}

	if count <= maxTokens {
		return log
	}

	ratio := float64(maxTokens) / float64(count)
	safeRatio := ratio * 0.90

	newLen := int(float64(len(log)) * safeRatio)

	if newLen == 0 {
		newLen = 1
	}

	logger.Info("Truncating log", "original_tokens", count, "limit", maxTokens, "original_len", len(log), "new_len", newLen)

	return log[:newLen] + "...[TRUNCATED]"
}
