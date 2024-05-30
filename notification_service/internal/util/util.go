package util

import (
	"github.com/aymerick/raymond"
	"github.com/gookit/slog"
)

func FormatMailMessage(data string, path string) string {
	reg, err := raymond.ParseFile("./templates/" + path)
	if err != nil {
		reg, err = raymond.ParseFile("../templates/" + path)
		if err != nil {
			slog.Fatalf("Failed to parse template: %v", err)
		}
	}

	ctx := map[string]interface{}{
		"data": data,
	}
	return reg.MustExec(ctx)
}
