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

	"github.com/goozt/gocargo/internal/auth"
	"github.com/goozt/gocargo/internal/build"
	"github.com/goozt/gocargo/internal/generate"
	"github.com/goozt/gocargo/internal/publish"
	"github.com/spf13/cobra"
)

var (
	flagName    string
	flagVersion string
	flagSummary string
	flagToken   string
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Build, package, and publish the Go binary as a Cargo crate",
	Long: `The publish command performs the full pipeline:

  1. Resolves a cargo registry token (flag → env → credentials file).
  2. Builds a static Go binary with CGO_ENABLED=0.
  3. Generates a minimal Rust shim crate that embeds the binary.
  4. Runs "cargo publish" to upload the crate to crates.io.`,
	RunE: runPublish,
}

func init() {
	publishCmd.Flags().StringVar(&flagName, "name", "", "crate package name (required)")
	publishCmd.Flags().StringVar(&flagVersion, "version", "", "semver version (required)")
	publishCmd.Flags().StringVar(&flagSummary, "summary", "", "short crate description (required)")
	publishCmd.Flags().StringVar(&flagToken, "token", "", "cargo registry token (optional)")

	_ = publishCmd.MarkFlagRequired("name")
	_ = publishCmd.MarkFlagRequired("version")
	_ = publishCmd.MarkFlagRequired("summary")

	rootCmd.AddCommand(publishCmd)
}

func runPublish(cmd *cobra.Command, args []string) error {
	// Clean previous build artifacts.
	if err := os.RemoveAll(".gocargo"); err != nil {
		return fmt.Errorf("clean .gocargo directory: %w", err)
	}

	// 1. Resolve authentication token.
	token := auth.ResolveToken(flagToken)

	// 2. Build the Go binary.
	fmt.Println("=> Building Go binary...")
	if err := build.GoBuild(); err != nil {
		return err
	}

	// 3. Generate Cargo.toml.
	fmt.Println("=> Generating Cargo.toml...")
	if err := generate.CargoToml(flagName, flagVersion, flagSummary); err != nil {
		return err
	}

	// 4. Generate src/main.rs.
	fmt.Println("=> Generating Rust shim (src/main.rs)...")
	if err := generate.MainRs(); err != nil {
		return err
	}

	// 5. Publish the crate.
	fmt.Println("=> Publishing crate...")
	if err := publish.CargoPublish(token); err != nil {
		return err
	}

	fmt.Printf("\nSuccess! Crate %q v%s published.\n", flagName, flagVersion)
	fmt.Printf("Users can now install it with: cargo install %s\n", flagName)
	return nil
}
