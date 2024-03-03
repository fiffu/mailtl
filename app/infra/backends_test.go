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
		backendContinue bool
		backendErr      error
	}{
		{
			backendContinue: true,
		},
		{
			backendContinue: false,
		},
		{
			backendErr: assert.AnError,
		},
	}
	for _, tc := range testCases {
		const task = backends.TaskSaveMail
		var expectRes = backends.NewResult(response.Canned.SuccessNoopCmd)
		desc := fmt.Sprintf(
			"when processor receives %s, and backend returns continue=%v err=%v",
			task, tc.backendContinue, tc.backendErr,
		)

		t.Run(desc, func(t *testing.T) {
			ctx := mock.Anything
			e := testfixtures.NewEnvelope()

			be := mocks.NewBackend(t)
			be.On("SaveMail", ctx, e).Return(tc.backendContinue, tc.backendErr)

			next := mocks.NewFixtureBackend(t)
			if tc.backendContinue {
				// Expect next processor to receive a call gracefully
				next.On("Process", e, task).
					Return(expectRes, nil)
			}

			constructor := MakeProcessor(be)
			proc := constructor()
			assert.NotNil(t, proc)

			res, err := proc(next).Process(e, task)
			if tc.backendErr != nil {
				expectRes = backends.NewResult(fmt.Sprintf("554 Error: %s", err))
				assert.Error(t, err)
			}
			assert.Equal(t, expectRes, res)
			assert.Equal(t, tc.backendErr, err, "processor should return the exact error from backend")
		})
	}
}

func TestMakeProcessor_UnknownTask(t *testing.T) {
	const unknownTask = backends.SelectTask(-1)

	e := testfixtures.NewEnvelope()
	be := mocks.NewBackend(t)

	next := mocks.NewFixtureBackend(t)
	proc := MakeProcessor(be)()

	res, err := proc(next).Process(e, unknownTask)
	assert.NoError(t, err)
	assert.Equal(t, backends.NewResult(response.Canned.SuccessNoopCmd), res)
}
