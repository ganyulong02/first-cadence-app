package workflows

/**

Workflows
Fault-oblivious stateful code is called workflow.
Used to chain together activities.
For more complex business logic, we can segregate this further into child workflows.

Activities
Workflow fault-oblivious code is immune to infrastructure failures. But we need to connect to other services where
failures are common. Hence codes that do all these are called activities.
All communication with external service should be using activities. Eg API call, Kafka Messages, Slack messages.
Think of it as the building blocks of Cadence.

Signals
Sometimes we want to have some influence while the workflow is running
Used to communicate with the workflow from external services such as Admin
 */

import (
	"context"
	"fmt"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldGroup"
const SignalName = "helloWorldSignal"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(helloworldActivity)
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func helloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "! How old are you!", nil
}

func Workflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("helloworld workflow started")
	var activityResult string
	err := workflow.ExecuteActivity(ctx, helloworldActivity, name).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	// After saying hello, the workflow will wait for you to inform it of your age!
	signalName := SignalName
	selector := workflow.NewSelector(ctx)
	var ageResult int

	for {
		signalChan := workflow.GetSignalChannel(ctx, signalName)
		selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
			c.Receive(ctx, &ageResult)
			workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
		})
		workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
		// Wait for signal
		selector.Select(ctx)

		// We can check the age and return an appropriate response
		if ageResult > 0 && ageResult < 150 {
			logger.Info("Workflow completed.", zap.String("NameResult", activityResult), zap.Int("AgeResult", ageResult))

			return fmt.Sprintf("Hello "+name+"! You are %v years old!", ageResult), nil
		} else {
			return "You can't be that old!", nil
		}
	}
}