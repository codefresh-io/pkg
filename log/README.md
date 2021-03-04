# Usage:
``` golang
package main

import (
	"context"

	"github.com/codefresh-io/foo/cmd/root"
	"github.com/codefresh-io/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
    lgrCfg := &log.LogrusConfig{Level: "info"}
	lgr := log.FromLogrus(logrus.NewEntry(logrus.StandardLogger()), lgrCfg)
	
    ctx = log.WithLogger(ctx, lgr)

	log.SetDefault(lgr) // if no context is provided to log.G() or the context does 
                        // not have a logger attached, you will get this logger.

	cmd := root.New(ctx)
	lgr.AddPFlags(cmd) // adds the logger flags to the command

	if err := cmd.Execute(); err != nil {
		log.G(ctx).Fatal(err)
	}
}

```