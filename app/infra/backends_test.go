package infra

import (
	"fmt"
	"testing"

	"github.com/fiffu/mailtl/app/infra/mocks"
	"github.com/fiffu/mailtl/testfixtures"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeProcessor_TaskSaveMail(t *testing.T) {
	testCases := []struct {
		stopProcessing bool
		backendErr     error
	}{
		{
			stopProcessing: true,
		},
		{
			stopProcessing: false,
		},
		{
			backendErr: assert.AnError,
		},
	}
	for _, tc := range testCases {
		const task = backends.TaskSaveMail
		desc := fmt.Sprintf(
			"when processor receives %s, and backend returns continue=%v err=%v",
			task, tc.stopProcessing, tc.backendErr != nil,
		)

		t.Run(desc, func(t *testing.T) {
			ctx := mock.Anything
			e := testfixtures.NewEnvelope()

			log := mocks.NewLogFacade(t)

			be := mocks.NewBackend(t)
			be.On("HandleTaskSaveMail", ctx, e).Return(tc.stopProcessing, tc.backendErr)
			be.On("Name").Maybe().Return("test_backend_name")

			next := mocks.NewFixtureBackend(t)

			var expectErr error
			expectRes := backends.NewResult(response.Canned.SuccessNoopCmd)
			switch {
			case tc.backendErr != nil:
				log.On("Errorf", ctx, tc.backendErr, "Processor %s errored", "test_backend_name").Return()
				expectErr = tc.backendErr
				expectRes = backends.NewResult(fmt.Sprintf("554 Error: %s", tc.backendErr))
			case tc.stopProcessing == false:
				next.On("Process", e, task).Return(expectRes, nil)
			}

			proc := MakeProcessor(log, be)()
			assert.NotNil(t, proc)

			res, err := proc(next).Process(e, task)
			assert.Equal(t, expectRes, res)
			assert.Equal(t, expectErr, err)
		})
	}
}

func TestMakeProcessor_UnknownTask(t *testing.T) {
	const unknownTask = backends.SelectTask(-1)

	e := testfixtures.NewEnvelope()
	be := mocks.NewBackend(t)
	log := mocks.NewLogFacade(t)

	next := mocks.NewFixtureBackend(t)
	next.On("Process", e, unknownTask).Return(nil, nil)

	proc := MakeProcessor(log, be)()

	res, err := proc(next).Process(e, unknownTask)
	assert.NoError(t, err)
	assert.Equal(t, nil, res)
}
