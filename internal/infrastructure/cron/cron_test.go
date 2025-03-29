package cron

import (
	"testing"
	"time"
)

func TestCronService_GetNextGenerationTime(t *testing.T) {
	service := NewCronService()

	tests := []struct {
		name       string
		expression string
		after      time.Time
		expected   time.Time
		wantErr    bool
	}{
		{
			name:       "valid cron expression - every minute",
			expression: "* * * * * *",
			after:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:   time.Date(2024, 1, 1, 12, 0, 1, 0, time.UTC),
			wantErr:    false,
		},
		{
			name:       "valid cron expression - every hour",
			expression: "0 0 * * * *",
			after:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			expected:   time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
			wantErr:    false,
		},
		{
			name:       "valid cron expression - daily at midnight",
			expression: "0 0 0 * * *",
			after:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:    false,
		},
		{
			name:       "invalid cron expression",
			expression: "invalid",
			after:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr:    true,
		},
		{
			name:       "empty cron expression",
			expression: "",
			after:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr:    true,
		},
		{
			name:       "cron expression with too few fields",
			expression: "* * *",
			after:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetNextGenerationTime(tt.expression, tt.after)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
