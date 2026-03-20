# gocargo

Use Cargo as a cross-platform package manager for Go projects. Embeds Go binaries into a Rust shim for seamless installation and execution.

## How It Works

`gocargo` takes your Go project, compiles it into a small static binary, and wraps it in a minimal Rust crate. When published to [crates.io](https://crates.io), anyone can install your Go tool with a single command:

```sh
cargo install your-tool-name
```

The generated Rust shim embeds the Go binary at compile time, extracts it to a temp file at runtime, and transparently proxies all arguments and I/O.

## Prerequisites

- **Go** (1.18+)
- **Rust / Cargo** (for publishing and for end-users installing your tool)

## Installation

```sh
go install github.com/goozt/gocargo@latest
```

Or build from source:

```sh
git clone https://github.com/goozt/gocargo.git
cd gocargo
go build -o gocargo .
```

## Usage

Run from the root of your Go project (where `go.mod` lives):

```sh
gocargo publish --name <crate-name> --version <semver> --summary "<description>"
```

### Example

```sh
gocargo publish --name gopher-tool --version 1.2.0 --summary "A high-performance Go CLI distributed via Cargo"
```

Then other users install it via:

```sh
cargo install gopher-tool
```

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Crate package name on crates.io |
| `--version` | Yes | Semantic version (e.g. `1.2.0`) |
| `--summary` | Yes | Short crate description |
| `--token` | No | Cargo registry token |

### Authentication

The registry token is resolved in this order:

1. `--token` flag
2. `GOCARGO_TOKEN` environment variable
3. `CARGO_REGISTRY_TOKEN` environment variable
4. `~/.cargo/credentials.toml` (automatic fallback)

## What Gets Generated

Running `gocargo publish` creates a `.gocargo/` directory:

```
.gocargo/
├── Cargo.toml        # Crate manifest with size-optimized release profile
├── bin/
│   └── app_bin       # Static Go binary (CGO_ENABLED=0, stripped)
└── src/
    └── main.rs       # Rust shim that embeds and executes the Go binary
```

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).
