// You can edit this code!
// Click here and start typing.
package main

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	svc, err := NewService(ctx, ServiceParams{
		AppPort: "8080",
	})

	if err != nil {
		logger.WithError(err).Fatal("Failed to create service")
	}

	<-svc.Run()

	logger.Info(ctx, "cmd: Exit")
}
