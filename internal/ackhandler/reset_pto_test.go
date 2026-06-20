package ackhandler

import (
	"testing"
	"time"

	"github.com/trzsz/quic-go/internal/monotime"
)

func TestResetPTO(t *testing.T) {
	// Create non-zero timestamps for testing
	now := monotime.Now()
	somePastTime := now.Add(-1 * time.Second)

	tests := []struct {
		name             string
		initialPtoCount  uint32
		initialAlarmTime monotime.Time
		now              monotime.Time
		expectPtoCount   uint32
		expectAlarmTime  monotime.Time
	}{
		{
			name:             "resets ptoCount and updates alarm time when alarm time is non-zero",
			initialPtoCount:  3,
			initialAlarmTime: somePastTime,
			now:              now,
			expectPtoCount:   0,
			expectAlarmTime:  now, // should be updated to the provided now
		},
		{
			name:             "only resets ptoCount and keeps alarm time zero when alarm time is zero",
			initialPtoCount:  2,
			initialAlarmTime: monotime.Time(0), // zero value
			now:              now,
			expectPtoCount:   0,
			expectAlarmTime:  monotime.Time(0), // should remain zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Initialize the minimal struct with the fields under test
			h := &sentPacketHandler{
				ptoCount: tt.initialPtoCount,
			}
			h.alarm.Time = tt.initialAlarmTime

			// 2. Call the method under test
			h.ResetPTO(tt.now)

			// 3. Verify the results match expectations
			if h.ptoCount != tt.expectPtoCount {
				t.Errorf("unexpected ptoCount: expected %d, got %d", tt.expectPtoCount, h.ptoCount)
			}
			if h.alarm.Time != tt.expectAlarmTime {
				t.Errorf("unexpected alarm.Time: expected %v, got %v", tt.expectAlarmTime, h.alarm.Time)
			}

			// 4. Verify that the loss detection timeout is not scheduled into the future
			if timeout := h.GetLossDetectionTimeout(); !timeout.IsZero() {
				if timeout.After(tt.now) {
					t.Errorf("unexpected loss detection timeout: scheduled after now (timeout: %v, now: %v)", timeout, tt.now)
				}
			}
		})
	}
}
