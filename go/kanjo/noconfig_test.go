package kanjobench

import "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

// Without this, pdfcpu writes config.yml and certificates into the user's
// configuration directory the first time a default configuration is built.
func init() { model.ConfigPath = "disable" }
