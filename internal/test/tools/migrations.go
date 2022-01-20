package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudquery/cloudquery/internal/test/provider"
	"github.com/cloudquery/cq-provider-sdk/migration"
)

func main() {
	if err := migration.Run(context.Background(), provider.Provider(), "internal/test/provider/migrations"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
