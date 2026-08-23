package tfdocsextras

import "embed"

// TemplateFS contains the structured per-section template used by both the
// tfdocs-extras CLI and the tfdocs-extras-plugin. The template breaks the
// output into reusable sections (inputs, outputs, objects, reference_links)
// that can be individually rendered and composed as needed.
//
//go:embed templates/structured.tmpl
var TemplateFS embed.FS
