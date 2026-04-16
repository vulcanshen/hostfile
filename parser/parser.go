package parser

import (
	"fmt"
	"net"
	"net/netip"
	"strings"
)

const (
	BlockStart = "#### hostfile >>>>>"
	BlockEnd   = "#### hostfile <<<<<"

	DisableIPPrefix     = "#[disable-ip]"
	DisableDomainPrefix = "#[disable-domain]"
)

type DisableType int

const (
	DisableNone   DisableType = iota
	DisableIP                 // entire line disabled
	DisableDomain             // single domain disabled
)

// ValidIP checks if a string is a valid IP address, including IPv6 with zone ID (e.g. fe80::1%lo0).
func ValidIP(s string) bool {
	if net.ParseIP(s) != nil {
		return true
	}
	// net.ParseIP doesn't handle zone IDs, try netip.ParseAddr
	_, err := netip.ParseAddr(s)
	return err == nil
}

// ValidDomain checks if a string is a valid domain name for hosts file use.
func ValidDomain(s string) bool {
	if len(s) == 0 || len(s) > 253 {
		return false
	}
	for _, label := range strings.Split(s, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
	}
	return true
}

// HostEntry represents a single entry in the managed block.
// For normal and disable-ip entries: IP + Domains.
// For disable-domain entries: IP + Domains (single domain that was disabled).
type HostEntry struct {
	IP          string
	Domains     []string
	DisableType DisableType
}

// ManagedBlock holds all entries within the managed block,
// plus any lines that failed to parse (stored as warnings).
type ManagedBlock struct {
	Entries  []HostEntry
	Warnings []string
}

// ParseLine parses a single line from within the managed block.
// Returns nil entry and nil error for empty/comment lines that should be skipped.
func ParseLine(line string) (*HostEntry, error) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, nil
	}

	// disable-domain: #[disable-domain] <domain> <ip>
	if strings.HasPrefix(trimmed, DisableDomainPrefix) {
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, DisableDomainPrefix))
		fields := strings.Fields(rest)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid disable-domain format: %s", line)
		}
		domain := fields[0]
		ip := fields[1]
		if !ValidIP(ip) {
			return nil, fmt.Errorf("invalid IP in disable-domain line: %s", ip)
		}
		return &HostEntry{
			IP:          ip,
			Domains:     []string{domain},
			DisableType: DisableDomain,
		}, nil
	}

	// disable-ip: #[disable-ip] <ip> <domain1> [domain2...]
	if strings.HasPrefix(trimmed, DisableIPPrefix) {
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, DisableIPPrefix))
		fields := strings.Fields(rest)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid disable-ip format: %s", line)
		}
		ip := fields[0]
		if !ValidIP(ip) {
			return nil, fmt.Errorf("invalid IP in disable-ip line: %s", ip)
		}
		return &HostEntry{
			IP:          ip,
			Domains:     fields[1:],
			DisableType: DisableIP,
		}, nil
	}

	// skip regular comments inside block
	if strings.HasPrefix(trimmed, "#") {
		return nil, nil
	}

	// normal line: <ip> <domain1> [domain2...]
	fields := strings.Fields(trimmed)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid hosts line (need IP + at least one domain): %s", line)
	}
	ip := fields[0]
	if !ValidIP(ip) {
		return nil, fmt.Errorf("invalid IP: %s", ip)
	}
	for _, d := range fields[1:] {
		if !ValidDomain(d) {
			return nil, fmt.Errorf("invalid domain %q (must be 1-253 chars, labels 1-63 chars)", d)
		}
	}
	return &HostEntry{
		IP:          ip,
		Domains:     fields[1:],
		DisableType: DisableNone,
	}, nil
}

// ParseBlock parses the content inside a managed block (without the marker lines).
func ParseBlock(content string) *ManagedBlock {
	block := &ManagedBlock{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		entry, err := ParseLine(line)
		if err != nil {
			block.Warnings = append(block.Warnings, fmt.Sprintf("warning: skipping line: %s (%v)", strings.TrimSpace(line), err))
			continue
		}
		if entry != nil {
			block.Entries = append(block.Entries, *entry)
		}
	}
	return block
}

// FormatEntry formats a single HostEntry back to a hosts file line.
func FormatEntry(entry *HostEntry) string {
	switch entry.DisableType {
	case DisableIP:
		return fmt.Sprintf("%s %s  %s", DisableIPPrefix, entry.IP, strings.Join(entry.Domains, " "))
	case DisableDomain:
		if len(entry.Domains) > 0 {
			return fmt.Sprintf("%s %s %s", DisableDomainPrefix, entry.Domains[0], entry.IP)
		}
		return ""
	default:
		return fmt.Sprintf("%s  %s", entry.IP, strings.Join(entry.Domains, " "))
	}
}

// FormatBlock formats a ManagedBlock back to string including the marker lines.
func FormatBlock(block *ManagedBlock) string {
	if block == nil || len(block.Entries) == 0 {
		return BlockStart + "\n" + BlockEnd
	}
	var lines []string
	lines = append(lines, BlockStart)
	for _, entry := range block.Entries {
		lines = append(lines, FormatEntry(&entry))
	}
	lines = append(lines, BlockEnd)
	return strings.Join(lines, "\n")
}
