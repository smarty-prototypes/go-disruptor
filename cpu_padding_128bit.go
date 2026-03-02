//go:build (arm64 && darwin) || ppc64 || ppc64le
// +build arm64,darwin ppc64 ppc64le

package disruptor

const CacheLineBytes = 128
