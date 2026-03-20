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

package build

import (
	"fmt"
	"os"
	"os/exec"
)

// GoBuild compiles the Go project in the current directory into a static binary
// at .gocargo/bin/app_bin using CGO_ENABLED=0 and stripped symbols.
func GoBuild() error {
	if err := os.MkdirAll(".gocargo/bin", 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", ".gocargo/bin/app_bin", ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	return nil
}
