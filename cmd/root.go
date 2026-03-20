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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocargo",
	Short: "Package Go binaries for distribution via cargo install",
	Long: `gocargo packages Go binaries into minimal Rust shim crates so that
end-users can install them with a single "cargo install" command.

The tool builds a static Go binary, generates a thin Rust wrapper that
embeds the binary, and publishes the resulting crate to crates.io.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
