package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Sync(t *testing.T) {
	t.Parallel()

	t.Run("changes dispatcher to sync mode", func(t *testing.T) {
		t.Parallel()

		// We dont need workers or a queue if we are running in sync mode
		dispatcher := NewDispatcher(0, 0)

		dispatcher.Sync()

		assert.True(t, dispatcher.syncMode)
	})

	t.Run("executes job synchronously without even queuing it", func(t *testing.T) {
		t.Parallel()

		// We dont need workers or a queue if we are running in sync mode
		dispatcher := NewDispatcher(0, 0)

		dispatcher.Sync()

		ran := false
		dispatcher.Run(func() {
			ran = true
		})

		assert.True(t, ran)
	})
}

func Test_Async(t *testing.T) {
	t.Parallel()

	t.Run("changes dispatcher to sync mode", func(t *testing.T) {
		t.Parallel()

		// We dont need workers or a queue if we are running in sync mode
		dispatcher := NewDispatcher(1, 1)

		dispatcher.Async()

		assert.False(t, dispatcher.syncMode)
	})
}
