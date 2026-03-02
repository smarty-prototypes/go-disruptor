//go:build 386 || amd64 || (arm64 && !darwin) || loong64 || riscv64 || wasm
// +build 386 amd64 arm64,!darwin loong64 riscv64 wasm

package disruptor

const CacheLineBytes = 64
