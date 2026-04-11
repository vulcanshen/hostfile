package parser

import (
	"strings"
	"testing"
)

func TestParseLine_Normal(t *testing.T) {
	entry, err := ParseLine("192.168.1.100  web.local api.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.IP != "192.168.1.100" {
		t.Errorf("expected IP 192.168.1.100, got %s", entry.IP)
	}
	if len(entry.Domains) != 2 || entry.Domains[0] != "web.local" || entry.Domains[1] != "api.local" {
		t.Errorf("expected [web.local api.local], got %v", entry.Domains)
	}
	if entry.DisableType != DisableNone {
		t.Errorf("expected DisableNone, got %v", entry.DisableType)
	}
}

func TestParseLine_IPv6(t *testing.T) {
	entry, err := ParseLine("::1  localhost6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.IP != "::1" {
		t.Errorf("expected IP ::1, got %s", entry.IP)
	}
	if len(entry.Domains) != 1 || entry.Domains[0] != "localhost6" {
		t.Errorf("expected [localhost6], got %v", entry.Domains)
	}
}

func TestParseLine_DisableIP(t *testing.T) {
	entry, err := ParseLine("#[disable-ip] 192.168.1.200  minio.company.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.IP != "192.168.1.200" {
		t.Errorf("expected IP 192.168.1.200, got %s", entry.IP)
	}
	if len(entry.Domains) != 1 || entry.Domains[0] != "minio.company.local" {
		t.Errorf("expected [minio.company.local], got %v", entry.Domains)
	}
	if entry.DisableType != DisableIP {
		t.Errorf("expected DisableIP, got %v", entry.DisableType)
	}
}

func TestParseLine_DisableDomain(t *testing.T) {
	entry, err := ParseLine("#[disable-domain] dockerhand.company.local 192.168.1.100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.IP != "192.168.1.100" {
		t.Errorf("expected IP 192.168.1.100, got %s", entry.IP)
	}
	if len(entry.Domains) != 1 || entry.Domains[0] != "dockerhand.company.local" {
		t.Errorf("expected [dockerhand.company.local], got %v", entry.Domains)
	}
	if entry.DisableType != DisableDomain {
		t.Errorf("expected DisableDomain, got %v", entry.DisableType)
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	entry, err := ParseLine("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Error("expected nil entry for empty line")
	}
}

func TestParseLine_Comment(t *testing.T) {
	entry, err := ParseLine("# this is a comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Error("expected nil entry for comment line")
	}
}

func TestParseLine_InvalidIP(t *testing.T) {
	_, err := ParseLine("not.an.ip web.local")
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}

func TestParseLine_NoDomainsNormal(t *testing.T) {
	_, err := ParseLine("192.168.1.1")
	if err == nil {
		t.Error("expected error for line with no domains")
	}
}

func TestParseLine_InvalidDisableIP(t *testing.T) {
	_, err := ParseLine("#[disable-ip] not-an-ip web.local")
	if err == nil {
		t.Error("expected error for invalid IP in disable-ip line")
	}
}

func TestParseLine_InvalidDisableDomain(t *testing.T) {
	_, err := ParseLine("#[disable-domain] only-one-field")
	if err == nil {
		t.Error("expected error for invalid disable-domain format")
	}
}

func TestParseBlock(t *testing.T) {
	content := `192.168.1.100  web.local api.local
#[disable-ip] 192.168.1.200  minio.company.local
#[disable-domain] docker.local 192.168.1.100
# a comment
bad line here
`
	block := ParseBlock(content)
	if len(block.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(block.Entries))
	}
	if len(block.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(block.Warnings), block.Warnings)
	}

	// entry 0: normal
	if block.Entries[0].IP != "192.168.1.100" || block.Entries[0].DisableType != DisableNone {
		t.Errorf("entry 0 mismatch: %+v", block.Entries[0])
	}
	// entry 1: disable-ip
	if block.Entries[1].IP != "192.168.1.200" || block.Entries[1].DisableType != DisableIP {
		t.Errorf("entry 1 mismatch: %+v", block.Entries[1])
	}
	// entry 2: disable-domain
	if block.Entries[2].IP != "192.168.1.100" || block.Entries[2].DisableType != DisableDomain {
		t.Errorf("entry 2 mismatch: %+v", block.Entries[2])
	}
}

func TestFormatEntry_Normal(t *testing.T) {
	entry := &HostEntry{IP: "192.168.1.100", Domains: []string{"web.local", "api.local"}, DisableType: DisableNone}
	got := FormatEntry(entry)
	expected := "192.168.1.100  web.local api.local"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatEntry_DisableIP(t *testing.T) {
	entry := &HostEntry{IP: "192.168.1.200", Domains: []string{"minio.local"}, DisableType: DisableIP}
	got := FormatEntry(entry)
	expected := "#[disable-ip] 192.168.1.200  minio.local"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatEntry_DisableDomain(t *testing.T) {
	entry := &HostEntry{IP: "192.168.1.100", Domains: []string{"docker.local"}, DisableType: DisableDomain}
	got := FormatEntry(entry)
	expected := "#[disable-domain] docker.local 192.168.1.100"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatBlock(t *testing.T) {
	block := &ManagedBlock{
		Entries: []HostEntry{
			{IP: "192.168.1.100", Domains: []string{"web.local"}, DisableType: DisableNone},
			{IP: "10.0.0.1", Domains: []string{"api.local"}, DisableType: DisableIP},
		},
	}
	got := FormatBlock(block)
	expected := "#### hostfile >>>>>\n192.168.1.100  web.local\n#[disable-ip] 10.0.0.1  api.local\n#### hostfile <<<<<"
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestFormatBlock_Empty(t *testing.T) {
	block := &ManagedBlock{}
	got := FormatBlock(block)
	expected := "#### hostfile >>>>>\n#### hostfile <<<<<"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatBlock_Nil(t *testing.T) {
	got := FormatBlock(nil)
	expected := "#### hostfile >>>>>\n#### hostfile <<<<<"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRoundTrip(t *testing.T) {
	original := `192.168.1.100  web.local api.local
#[disable-ip] 192.168.1.200  minio.local
#[disable-domain] docker.local 10.0.0.1`

	block := ParseBlock(original)
	if len(block.Warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", block.Warnings)
	}

	formatted := FormatBlock(block)
	// re-parse the formatted output (strip markers)
	lines := []string{}
	for _, line := range splitLines(formatted) {
		if line != BlockStart && line != BlockEnd {
			lines = append(lines, line)
		}
	}
	reparsed := ParseBlock(joinLines(lines))
	if len(reparsed.Entries) != len(block.Entries) {
		t.Fatalf("round trip entry count mismatch: %d vs %d", len(reparsed.Entries), len(block.Entries))
	}
	for i, e := range block.Entries {
		r := reparsed.Entries[i]
		if e.IP != r.IP || e.DisableType != r.DisableType || len(e.Domains) != len(r.Domains) {
			t.Errorf("round trip mismatch at %d: %+v vs %+v", i, e, r)
		}
	}
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
