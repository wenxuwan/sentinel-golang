package base

import (
	"fmt"
	"time"
)

type BlockType uint8

const (
	BlockTypeUnknown BlockType = iota
	BlockTypeFlow
	BlockTypeIsolation
	BlockTypeCircuitBreaking
	BlockTypeSystemFlow
	BlockTypeHotSpotParamFlow
)

func (t BlockType) String() string {
	switch t {
	case BlockTypeUnknown:
		return "Unknown"
	case BlockTypeFlow:
		return "FlowControl"
	case BlockTypeIsolation:
		return "BlockTypeIsolation"
	case BlockTypeCircuitBreaking:
		return "CircuitBreaking"
	case BlockTypeSystemFlow:
		return "System"
	case BlockTypeHotSpotParamFlow:
		return "HotSpotParamFlow"
	default:
		return fmt.Sprintf("%d", t)
	}
}

type TokenResultStatus uint8

const (
	ResultStatusPass TokenResultStatus = iota
	ResultStatusBlocked
	ResultStatusShouldWait
)

func (s TokenResultStatus) String() string {
	switch s {
	case ResultStatusPass:
		return "ResultStatusPass"
	case ResultStatusBlocked:
		return "ResultStatusBlocked"
	case ResultStatusShouldWait:
		return "ResultStatusShouldWait"
	default:
		return "Undefined"
	}
}

type TokenResult struct {
	status TokenResultStatus

	blockErr    *BlockError
	nanosToWait time.Duration
}

func (r *TokenResult) DeepCopyFrom(newResult *TokenResult) {
	r.status = newResult.status
	r.nanosToWait = newResult.nanosToWait
	if r.blockErr == nil {
		r.blockErr = &BlockError{
			blockType:     newResult.blockErr.blockType,
			blockMsg:      newResult.blockErr.blockMsg,
			rule:          newResult.blockErr.rule,
			snapshotValue: newResult.blockErr.snapshotValue,
		}
	} else {
		// TODO: review the reusing logic
		r.blockErr.blockType = newResult.blockErr.blockType
		r.blockErr.blockMsg = newResult.blockErr.blockMsg
		r.blockErr.rule = newResult.blockErr.rule
		r.blockErr.snapshotValue = newResult.blockErr.snapshotValue
	}
}

func (r *TokenResult) ResetToPass() {
	r.status = ResultStatusPass
	r.blockErr = nil
	r.nanosToWait = 0
}

func (r *TokenResult) ResetToBlocked(blockType BlockType) {
	r.status = ResultStatusBlocked
	if r.blockErr == nil {
		r.blockErr = NewBlockError(blockType)
	} else {
		r.blockErr.blockType = blockType
		r.blockErr.blockMsg = ""
		r.blockErr.rule = nil
		r.blockErr.snapshotValue = nil
	}
	r.nanosToWait = 0
}

func (r *TokenResult) ResetToBlockedWithMessage(blockType BlockType, blockMsg string) {
	r.status = ResultStatusBlocked
	if r.blockErr == nil {
		r.blockErr = NewBlockErrorWithMessage(blockType, blockMsg)
	} else {
		r.blockErr.blockType = blockType
		r.blockErr.blockMsg = blockMsg
		r.blockErr.rule = nil
		r.blockErr.snapshotValue = nil
	}
	r.nanosToWait = 0
}

func (r *TokenResult) ResetToBlockedWithCause(blockType BlockType, blockMsg string, rule SentinelRule, snapshot interface{}) {
	r.status = ResultStatusBlocked
	if r.blockErr == nil {
		r.blockErr = NewBlockErrorWithCause(blockType, blockMsg, rule, snapshot)
	} else {
		r.blockErr.blockType = blockType
		r.blockErr.blockMsg = blockMsg
		r.blockErr.rule = rule
		r.blockErr.snapshotValue = snapshot
	}
	r.nanosToWait = 0
}

func (r *TokenResult) IsPass() bool {
	return r.status == ResultStatusPass
}

func (r *TokenResult) IsBlocked() bool {
	return r.status == ResultStatusBlocked
}

func (r *TokenResult) Status() TokenResultStatus {
	return r.status
}

func (r *TokenResult) BlockError() *BlockError {
	return r.blockErr
}

func (r *TokenResult) NanosToWait() time.Duration {
	return r.nanosToWait
}

func (r *TokenResult) String() string {
	var blockMsg string
	if r.blockErr == nil {
		blockMsg = "none"
	} else {
		blockMsg = r.blockErr.Error()
	}
	return fmt.Sprintf("TokenResult{status=%s, blockErr=%s, nanosToWait=%d}", r.status.String(), blockMsg, r.nanosToWait)
}

func NewTokenResultPass() *TokenResult {
	return &TokenResult{
		status:      ResultStatusPass,
		blockErr:    nil,
		nanosToWait: 0,
	}
}

func NewTokenResultBlocked(blockType BlockType) *TokenResult {
	return &TokenResult{
		status:      ResultStatusBlocked,
		blockErr:    NewBlockError(blockType),
		nanosToWait: 0,
	}
}

func NewTokenResultBlockedWithMessage(blockType BlockType, blockMsg string) *TokenResult {
	return &TokenResult{
		status:      ResultStatusBlocked,
		blockErr:    NewBlockErrorWithMessage(blockType, blockMsg),
		nanosToWait: 0,
	}
}

func NewTokenResultBlockedWithCause(blockType BlockType, blockMsg string, rule SentinelRule, snapshot interface{}) *TokenResult {
	return &TokenResult{
		status:      ResultStatusBlocked,
		blockErr:    NewBlockErrorWithCause(blockType, blockMsg, rule, snapshot),
		nanosToWait: 0,
	}
}

func NewTokenResultShouldWait(waitNs time.Duration) *TokenResult {
	return &TokenResult{
		status:      ResultStatusShouldWait,
		blockErr:    nil,
		nanosToWait: waitNs,
	}
}
