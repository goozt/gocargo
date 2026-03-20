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

package auth

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// credentials represents the structure of ~/.cargo/credentials.toml.
type credentials struct {
	Registry   registryEntry            `toml:"registry"`
	Registries map[string]registryEntry `toml:"registries"`
}

type registryEntry struct {
	Token string `toml:"token"`
}

// ResolveToken returns a cargo registry token using the following priority:
//  1. The explicit flag value
//  2. GOCARGO_TOKEN environment variable
//  3. CARGO_REGISTRY_TOKEN environment variable
//  4. Token from ~/.cargo/credentials.toml
//
// Returns an empty string if no token is found.
func ResolveToken(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}

	if t := os.Getenv("GOCARGO_TOKEN"); t != "" {
		return t
	}

	if t := os.Getenv("CARGO_REGISTRY_TOKEN"); t != "" {
		return t
	}

	return readCredentialsToken()
}

func readCredentialsToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Cargo supports both credentials.toml and credentials (no extension).
	for _, name := range []string{"credentials.toml", "credentials"} {
		path := filepath.Join(home, ".cargo", name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var creds credentials
		if err := toml.Unmarshal(data, &creds); err != nil {
			continue
		}

		if creds.Registry.Token != "" {
			return creds.Registry.Token
		}

		if entry, ok := creds.Registries["crates-io"]; ok && entry.Token != "" {
			return entry.Token
		}
	}

	return ""
}
