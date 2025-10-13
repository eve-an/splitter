package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/feature"
)

func ListAllFeatures(w http.ResponseWriter, r *http.Request) {
	features := []feature.Feature{
		{
			Name:         "red_button",
			Descritption: "Test wether the new red butto is better than the current green one.",
			Active:       true,
			Variants: []feature.Variant{
				{
					Name:   "red_button",
					Weight: 50,
				},
				{
					Name:   "green_button",
					Weight: 50,
				},
			},
		},
	}

	if err := json.NewEncoder(w).Encode(features); err != nil {
		slog.Error("failed marshalling feature", slog.Any("error", err))
	}
}
