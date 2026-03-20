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

package publish

import (
	"fmt"
	"os"
	"os/exec"
)

// CargoPublish runs `cargo publish` for the generated crate.
// If token is non-empty it is passed via --token.
func CargoPublish(token string) error {
	if _, err := exec.LookPath("cargo"); err != nil {
		return fmt.Errorf("cargo not found in PATH: %w", err)
	}

	args := []string{"publish", "--manifest-path", ".gocargo/Cargo.toml"}
	if token != "" {
		args = append(args, "--token", token)
	}

	cmd := exec.Command("cargo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cargo publish: %w", err)
	}

	return nil
}
