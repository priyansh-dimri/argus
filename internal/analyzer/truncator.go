package analyzer

import (
	"context"

	"github.com/priyansh-dimri/argus/pkg/logger"
)

func TruncateLog(ctx context.Context, client AIClient, log string, maxTokens int) string {
	logger.Info("Starting log truncation check",
		"component", "truncator",
		"log_length", len(log),
		"max_tokens", maxTokens,
	)

	if len(log) < maxTokens {
		logger.Info("Log length is below max tokens, no truncation needed",
			"component", "truncator",
			"log_length", len(log),
			"max_tokens", maxTokens,
		)
		return log
	}

	count, err := client.CountTokens(ctx, log)
	if err != nil {
		logger.Warn("Failed to count tokens during truncation, falling back to safe mode",
			"component", "truncator",
			"error", err,
			"log_length", len(log),
		)

		safeCharLimit := maxTokens * 3
		if len(log) > safeCharLimit {
			logger.Info("Applying safe mode truncation",
				"component", "truncator",
				"original_length", len(log),
				"safe_char_limit", safeCharLimit,
			)
			return log[:safeCharLimit] + "...[TRUNCATED_SAFE_MODE]"
		}

		logger.Info("Log is within safe character limit, no truncation needed",
			"component", "truncator",
			"log_length", len(log),
			"safe_char_limit", safeCharLimit,
		)
		return log
	}

	logger.Info("Token count completed successfully",
		"component", "truncator",
		"token_count", count,
		"max_tokens", maxTokens,
	)

	if count <= maxTokens {
		logger.Info("Token count is within limit, no truncation needed",
			"component", "truncator",
			"token_count", count,
			"max_tokens", maxTokens,
		)
		return log
	}

	ratio := float64(maxTokens) / float64(count)
	safeRatio := ratio * 0.90

	newLen := int(float64(len(log)) * safeRatio)

	if newLen == 0 {
		newLen = 1
		logger.Warn("Calculated truncation length was 0, setting to 1",
			"component", "truncator",
		)
	}

	logger.Info("Truncating log to fit token limit",
		"component", "truncator",
		"original_tokens", count,
		"max_tokens", maxTokens,
		"original_length", len(log),
		"new_length", newLen,
		"reduction_ratio", safeRatio,
	)

	return log[:newLen] + "...[TRUNCATED]"
}
