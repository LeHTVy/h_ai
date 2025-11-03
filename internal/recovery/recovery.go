package recovery

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type RecoveryAction string

const (
	RetryWithBackoff          RecoveryAction = "retry_with_backoff"
	RetryWithReducedScope     RecoveryAction = "retry_with_reduced_scope"
	SwitchToAlternativeTool   RecoveryAction = "switch_to_alternative_tool"
	AdjustParameters          RecoveryAction = "adjust_parameters"
	EscalateToHuman           RecoveryAction = "escalate_to_human"
	GracefulDegradation       RecoveryAction = "graceful_degradation"
	AbortOperation            RecoveryAction = "abort_operation"
)

type RecoveryStrategy struct {
	Action           RecoveryAction
	Parameters       map[string]interface{}
	BackoffMultiplier float64
	MaxAttempts      int
}

type RecoveryManager struct {
	logger   *zap.Logger
	strategies map[string]RecoveryStrategy
}

func New(logger *zap.Logger) *RecoveryManager {
	rm := &RecoveryManager{
		logger:    logger,
		strategies: make(map[string]RecoveryStrategy),
	}

	// Initialize default strategies
	rm.initDefaultStrategies()
	return rm
}

func (rm *RecoveryManager) initDefaultStrategies() {
	// Default strategy for network tools
	rm.strategies["nmap"] = RecoveryStrategy{
		Action:            RetryWithBackoff,
		BackoffMultiplier: 2.0,
		MaxAttempts:       3,
		Parameters: map[string]interface{}{
			"initial_delay": 5,
			"max_delay":     60,
		},
	}

	// Default strategy for web tools
	rm.strategies["gobuster"] = RecoveryStrategy{
		Action:            RetryWithBackoff,
		BackoffMultiplier: 1.5,
		MaxAttempts:       2,
		Parameters: map[string]interface{}{
			"initial_delay": 3,
			"max_delay":     30,
		},
	}
}

func (rm *RecoveryManager) GetStrategy(toolName string, errorType string) RecoveryStrategy {
	// Try to get tool-specific strategy
	if strategy, found := rm.strategies[toolName]; found {
		return strategy
	}

	// Default strategy
	return RecoveryStrategy{
		Action:            RetryWithBackoff,
		BackoffMultiplier: 2.0,
		MaxAttempts:       3,
		Parameters: map[string]interface{}{
			"initial_delay": 5,
			"max_delay":     60,
		},
	}
}

func (rm *RecoveryManager) CalculateBackoff(attempt int, strategy RecoveryStrategy) time.Duration {
	initialDelay := 5 * time.Second
	if delay, ok := strategy.Parameters["initial_delay"].(int); ok {
		initialDelay = time.Duration(delay) * time.Second
	}

	maxDelay := 60 * time.Second
	if delay, ok := strategy.Parameters["max_delay"].(int); ok {
		maxDelay = time.Duration(delay) * time.Second
	}

	backoff := float64(initialDelay) * (strategy.BackoffMultiplier * float64(attempt-1))
	if backoff > float64(maxDelay) {
		backoff = float64(maxDelay)
	}

	return time.Duration(backoff)
}

func (rm *RecoveryManager) HandleError(toolName string, errorMsg string, attempt int) (RecoveryStrategy, error) {
	strategy := rm.GetStrategy(toolName, "")

	if attempt >= strategy.MaxAttempts {
		return RecoveryStrategy{
			Action: EscalateToHuman,
		}, fmt.Errorf("max attempts reached for %s", toolName)
	}

	rm.logger.Info("Recovery strategy selected",
		zap.String("tool", toolName),
		zap.String("action", string(strategy.Action)),
		zap.Int("attempt", attempt))

	return strategy, nil
}
