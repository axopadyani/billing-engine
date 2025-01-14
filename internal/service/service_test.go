package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/axopadyani/billing-engine/internal/test/mock/repository"
)

var testTimeout = 10 * time.Second

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)

	svc := NewService(mockRepo)
	if svc == nil {
		t.Error("expecting service to be created")
	}
}
