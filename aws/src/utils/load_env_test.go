package utils_test

import (
	"context"
	"testing"

	"github.com/bengosborn/roomiez/aws/utils"
)

func TestLoadEnv(t *testing.T) {
	ctx := context.Background()

	t.Run("Load env", func(t *testing.T) {
		if _, err := utils.LoadEnv(ctx); err != nil {
			t.Error(err)
		}
	})
}
