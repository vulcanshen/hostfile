package manager

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/vulcanshen/hostfile/parser"
)

// ReadHostsFile reads a hosts file and splits it into three parts:
// content before the managed block, the managed block itself, and content after.
// If no managed block exists, block is empty and all content goes to "after"
// (so that a new block will be prepended at the top).
func ReadHostsFile(path string) (before string, block *parser.ManagedBlock, after string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", &parser.ManagedBlock{}, "", nil
		}
		return "", nil, "", err
	}

	content := string(data)
	startIdx := strings.Index(content, parser.BlockStart)
	endIdx := strings.Index(content, parser.BlockEnd)

	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		// no valid block found — all content goes to "after" so new block is prepended
		return "", &parser.ManagedBlock{}, content, nil
	}

	before = content[:startIdx]
	blockContent := content[startIdx+len(parser.BlockStart) : endIdx]
	after = content[endIdx+len(parser.BlockEnd):]

	// trim leading newline from block content
	blockContent = strings.TrimPrefix(blockContent, "\n")
	blockContent = strings.TrimSuffix(blockContent, "\n")

	block = parser.ParseBlock(blockContent)
	for _, w := range block.Warnings {
		fmt.Fprintln(os.Stderr, w)
	}

	return before, block, after, nil
}

// WriteHostsFile writes back the hosts file with the managed block.
func WriteHostsFile(path string, before string, block *parser.ManagedBlock, after string) error {
	formatted := parser.FormatBlock(block)

	var content string
	if before == "" && after == "" {
		content = formatted + "\n"
	} else if before == "" {
		// block at top, then existing content
		content = formatted + "\n" + after
	} else {
		content = before + formatted + "\n" + after
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// isIP returns true if the string is a valid IP address.
func isIP(s string) bool {
	return net.ParseIP(s) != nil
}

// Add adds domains to an IP in the managed block.
// If the IP already exists, domains are appended.
// Returns a ConflictInfo if any domain is already mapped to a different IP.
type ConflictInfo struct {
	Domain    string
	CurrentIP string
}

func Add(block *parser.ManagedBlock, ip string, domains []string) ([]ConflictInfo, error) {
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	var conflicts []ConflictInfo
	for _, domain := range domains {
		for _, entry := range block.Entries {
			if entry.IP == ip {
				continue
			}
			if entry.DisableType == parser.DisableDomain {
				if len(entry.Domains) > 0 && entry.Domains[0] == domain {
					conflicts = append(conflicts, ConflictInfo{Domain: domain, CurrentIP: entry.IP})
				}
				continue
			}
			for _, d := range entry.Domains {
				if d == domain {
					conflicts = append(conflicts, ConflictInfo{Domain: domain, CurrentIP: entry.IP})
				}
			}
		}
	}

	if len(conflicts) > 0 {
		return conflicts, nil
	}

	// find existing entry for this IP
	for i, entry := range block.Entries {
		if entry.IP == ip && entry.DisableType == parser.DisableNone {
			for _, domain := range domains {
				if !containsDomain(entry.Domains, domain) {
					block.Entries[i].Domains = append(block.Entries[i].Domains, domain)
				}
			}
			return nil, nil
		}
	}

	// new entry
	block.Entries = append(block.Entries, parser.HostEntry{
		IP:          ip,
		Domains:     domains,
		DisableType: parser.DisableNone,
	})
	return nil, nil
}

// AddForce adds domains to an IP, moving them from any conflicting IP.
func AddForce(block *parser.ManagedBlock, ip string, domains []string) error {
	for _, domain := range domains {
		removeDomainFromBlock(block, domain)
	}
	_, err := Add(block, ip, domains)
	return err
}

// RemoveResult describes what happened after a remove operation.
type RemoveResult struct {
	// LastDomain is true if the target was a domain and it was the last one on its IP line.
	// The caller should confirm with the user whether to remove the entire IP line.
	LastDomain    bool
	LastDomainIP  string
	Removed       bool
}

// Remove removes an IP or domain from the managed block.
// If removing a domain that is the last one on an IP line, it returns LastDomain=true
// so the caller can ask the user for confirmation.
func Remove(block *parser.ManagedBlock, target string) RemoveResult {
	if isIP(target) {
		return removeIP(block, target)
	}
	return removeDomain(block, target)
}

// RemoveConfirmed removes the IP line after user confirmed removing the last domain.
func RemoveConfirmed(block *parser.ManagedBlock, ip string) {
	removeIP(block, ip)
}

func removeIP(block *parser.ManagedBlock, ip string) RemoveResult {
	newEntries := make([]parser.HostEntry, 0, len(block.Entries))
	removed := false
	for _, entry := range block.Entries {
		if entry.IP == ip {
			removed = true
			continue
		}
		newEntries = append(newEntries, entry)
	}
	block.Entries = newEntries
	return RemoveResult{Removed: removed}
}

func removeDomain(block *parser.ManagedBlock, domain string) RemoveResult {
	for i, entry := range block.Entries {
		if entry.DisableType == parser.DisableDomain {
			if len(entry.Domains) > 0 && entry.Domains[0] == domain {
				// remove this disable-domain entry
				block.Entries = append(block.Entries[:i], block.Entries[i+1:]...)
				return RemoveResult{Removed: true}
			}
			continue
		}

		idx := domainIndex(entry.Domains, domain)
		if idx == -1 {
			continue
		}

		if len(entry.Domains) == 1 {
			return RemoveResult{LastDomain: true, LastDomainIP: entry.IP}
		}

		block.Entries[i].Domains = append(entry.Domains[:idx], entry.Domains[idx+1:]...)
		return RemoveResult{Removed: true}
	}
	return RemoveResult{}
}

// Search searches the managed block for entries matching the query (IP or domain).
func Search(block *parser.ManagedBlock, query string) []parser.HostEntry {
	var results []parser.HostEntry
	if isIP(query) {
		for _, entry := range block.Entries {
			if entry.IP == query {
				results = append(results, entry)
			}
		}
	} else {
		for _, entry := range block.Entries {
			if entry.DisableType == parser.DisableDomain {
				if len(entry.Domains) > 0 && entry.Domains[0] == query {
					results = append(results, entry)
				}
				continue
			}
			for _, d := range entry.Domains {
				if d == query {
					results = append(results, entry)
					break
				}
			}
		}
	}
	return results
}

// List returns all entries in the managed block.
func List(block *parser.ManagedBlock) []parser.HostEntry {
	return block.Entries
}

// Enable enables a disabled entry by IP or domain.
// enable <domain>: removes the disable-domain entry and adds the domain back to its IP line.
// enable <ip>: changes disable-ip entries for that IP to normal entries.
func Enable(block *parser.ManagedBlock, target string) error {
	if isIP(target) {
		return enableIP(block, target)
	}
	return enableDomain(block, target)
}

func enableIP(block *parser.ManagedBlock, ip string) error {
	found := false
	for i, entry := range block.Entries {
		if entry.IP == ip && entry.DisableType == parser.DisableIP {
			block.Entries[i].DisableType = parser.DisableNone
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no disabled entry found for IP %s", ip)
	}
	return nil
}

func enableDomain(block *parser.ManagedBlock, domain string) error {
	// find the disable-domain entry
	for i, entry := range block.Entries {
		if entry.DisableType == parser.DisableDomain && len(entry.Domains) > 0 && entry.Domains[0] == domain {
			ip := entry.IP
			// remove the disable-domain entry
			block.Entries = append(block.Entries[:i], block.Entries[i+1:]...)
			// add domain back to the IP line (or create one)
			for j, e := range block.Entries {
				if e.IP == ip && (e.DisableType == parser.DisableNone || e.DisableType == parser.DisableIP) {
					block.Entries[j].Domains = append(block.Entries[j].Domains, domain)
					return nil
				}
			}
			// no existing line for this IP, create one
			block.Entries = append(block.Entries, parser.HostEntry{
				IP:          ip,
				Domains:     []string{domain},
				DisableType: parser.DisableNone,
			})
			return nil
		}
	}
	return fmt.Errorf("no disabled entry found for domain %s", domain)
}

// DisableResult describes what happened after a disable operation.
type DisableResult struct {
	// LastDomain is true if the target domain is the last active one on its IP line.
	// The caller should confirm with the user whether to disable the entire IP.
	LastDomain   bool
	LastDomainIP string
	Done         bool
}

// Disable disables an entry by IP or domain.
func Disable(block *parser.ManagedBlock, target string) (DisableResult, error) {
	if isIP(target) {
		return disableIP(block, target)
	}
	return disableDomain(block, target)
}

// DisableIPConfirmed disables the entire IP line after user confirmed.
func DisableIPConfirmed(block *parser.ManagedBlock, ip string) {
	disableIP(block, ip)
}

func disableIP(block *parser.ManagedBlock, ip string) (DisableResult, error) {
	found := false
	for i, entry := range block.Entries {
		if entry.IP == ip && entry.DisableType == parser.DisableNone {
			block.Entries[i].DisableType = parser.DisableIP
			found = true
		}
	}
	if !found {
		return DisableResult{}, fmt.Errorf("no active entry found for IP %s", ip)
	}
	return DisableResult{Done: true}, nil
}

func disableDomain(block *parser.ManagedBlock, domain string) (DisableResult, error) {
	for i, entry := range block.Entries {
		if entry.DisableType != parser.DisableNone {
			continue
		}
		idx := domainIndex(entry.Domains, domain)
		if idx == -1 {
			continue
		}

		if len(entry.Domains) == 1 {
			return DisableResult{LastDomain: true, LastDomainIP: entry.IP}, nil
		}

		ip := entry.IP
		// remove domain from the entry
		block.Entries[i].Domains = append(entry.Domains[:idx], entry.Domains[idx+1:]...)
		// add a disable-domain entry
		block.Entries = append(block.Entries, parser.HostEntry{
			IP:          ip,
			Domains:     []string{domain},
			DisableType: parser.DisableDomain,
		})
		return DisableResult{Done: true}, nil
	}
	return DisableResult{}, fmt.Errorf("no active entry found for domain %s", domain)
}

// Clean removes all entries from the managed block.
func Clean(block *parser.ManagedBlock) {
	block.Entries = nil
}

// Apply replaces the managed block with content parsed from the given string.
func Apply(block *parser.ManagedBlock, content string) {
	newBlock := parser.ParseBlock(content)
	for _, w := range newBlock.Warnings {
		fmt.Fprintln(os.Stderr, w)
	}
	block.Entries = newBlock.Entries
}

// Merge merges content into the existing managed block.
// Entries with the same IP are merged (domains appended).
// New IPs are added as new entries.
func Merge(block *parser.ManagedBlock, content string) {
	newBlock := parser.ParseBlock(content)
	for _, w := range newBlock.Warnings {
		fmt.Fprintln(os.Stderr, w)
	}

	for _, newEntry := range newBlock.Entries {
		if newEntry.DisableType != parser.DisableNone {
			block.Entries = append(block.Entries, newEntry)
			continue
		}
		merged := false
		for i, existing := range block.Entries {
			if existing.IP == newEntry.IP && existing.DisableType == parser.DisableNone {
				for _, d := range newEntry.Domains {
					if !containsDomain(existing.Domains, d) {
						block.Entries[i].Domains = append(block.Entries[i].Domains, d)
					}
				}
				merged = true
				break
			}
		}
		if !merged {
			block.Entries = append(block.Entries, newEntry)
		}
	}
}

// helper functions

func containsDomain(domains []string, domain string) bool {
	for _, d := range domains {
		if d == domain {
			return true
		}
	}
	return false
}

func domainIndex(domains []string, domain string) int {
	for i, d := range domains {
		if d == domain {
			return i
		}
	}
	return -1
}

func removeDomainFromBlock(block *parser.ManagedBlock, domain string) {
	for i := len(block.Entries) - 1; i >= 0; i-- {
		entry := &block.Entries[i]
		if entry.DisableType == parser.DisableDomain {
			if len(entry.Domains) > 0 && entry.Domains[0] == domain {
				block.Entries = append(block.Entries[:i], block.Entries[i+1:]...)
			}
			continue
		}
		idx := domainIndex(entry.Domains, domain)
		if idx != -1 {
			entry.Domains = append(entry.Domains[:idx], entry.Domains[idx+1:]...)
			if len(entry.Domains) == 0 {
				block.Entries = append(block.Entries[:i], block.Entries[i+1:]...)
			}
		}
	}
}
