package exporter

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

func validatePromLabelName(name string) error {
	if name == "" {
		return fmt.Errorf("empty label name")
	}
	if strings.HasPrefix(name, model.ReservedLabelPrefix) {
		return fmt.Errorf("label name %q uses reserved prefix %q", name, model.ReservedLabelPrefix)
	}
	if !model.LegacyValidation.IsValidLabelName(name) {
		return fmt.Errorf("invalid label name %q", name)
	}
	return nil
}

func validatePromMetricName(name string) error {
	if name == "" {
		return fmt.Errorf("empty metric name")
	}
	if !model.LegacyValidation.IsValidMetricName(name) {
		return fmt.Errorf("invalid metric name %q", name)
	}
	return nil
}
