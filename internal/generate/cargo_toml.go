// gocargo - Package Go binaries into Rust shim crates for cargo install
// Copyright (C) 2026 gocargo contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package generate

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

const cargoTomlTmpl = `[package]
name = "{{.Name}}"
version = "{{.Version}}"
edition = "2021"
description = "{{.Summary}}"
license = "GPL-3.0"

[profile.release]
strip = true
opt-level = "z"
panic = "abort"
`

type cargoTomlData struct {
	Name    string
	Version string
	Summary string
}

// CargoToml generates the .gocargo/Cargo.toml file with the given package metadata.
func CargoToml(name, version, summary string) error {
	t, err := template.New("Cargo.toml").Parse(cargoTomlTmpl)
	if err != nil {
		return fmt.Errorf("parse Cargo.toml template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, cargoTomlData{
		Name:    name,
		Version: version,
		Summary: summary,
	}); err != nil {
		return fmt.Errorf("execute Cargo.toml template: %w", err)
	}

	if err := os.WriteFile(".gocargo/Cargo.toml", buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write Cargo.toml: %w", err)
	}

	return nil
}
