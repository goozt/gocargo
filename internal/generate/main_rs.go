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
	"fmt"
	"os"
)

const mainRsBody = `
use std::io::Write;

fn main() {
    let bytes = include_bytes!("../bin/app_bin");

    let mut path = std::env::temp_dir();
    let unique = format!(
        "gocargo-{}-{}",
        std::process::id(),
        std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos()
    );

    #[cfg(unix)]
    path.push(&unique);
    #[cfg(windows)]
    path.push(format!("{}.exe", unique));

    let mut f = std::fs::File::create(&path).expect("failed to create temp file");
    f.write_all(bytes).expect("failed to write binary");
    drop(f);

    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        std::fs::set_permissions(&path, std::fs::Permissions::from_mode(0o755))
            .expect("failed to set permissions");
    }

    let args: Vec<String> = std::env::args().skip(1).collect();
    let status = std::process::Command::new(&path)
        .args(&args)
        .stdin(std::process::Stdio::inherit())
        .stdout(std::process::Stdio::inherit())
        .stderr(std::process::Stdio::inherit())
        .status()
        .expect("failed to execute binary");

    let _ = std::fs::remove_file(&path);
    std::process::exit(status.code().unwrap_or(1));
}
`

// MainRs generates the .gocargo/src/main.rs Rust shim that embeds and executes
// the Go binary.
func MainRs() error {
	if err := os.MkdirAll(".gocargo/src", 0o755); err != nil {
		return fmt.Errorf("create src directory: %w", err)
	}

	content := mainRsBody

	if err := os.WriteFile(".gocargo/src/main.rs", []byte(content), 0o644); err != nil {
		return fmt.Errorf("write main.rs: %w", err)
	}

	return nil
}
