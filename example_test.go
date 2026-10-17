package surge_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-surge"
)

// This is the README example. It needs a real token, so it only compiles.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := surge.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization", "Bearer "+os.Getenv("SURGE_API_KEY"),
			)

			return nil
		},
	)

	client, err := surge.NewClientWithResponses("https://api.surge.app", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListAccountsWithResponse(
		context.Background(),
		&surge.ListAccountsParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
