package casbin_middleware

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPathMappingFromMenuFile(t *testing.T) {
	tempDir := t.TempDir()
	menuFile := filepath.Join(tempDir, "menu.yaml")
	content := `menu:
  - id: parent
    children:
      - id: account_management
        api_paths:
          - /api/v1/identity/users
          - /api/v1/identity/users/*
      - id: role_permissions
        api_paths:
          - /api/v1/permission/roles
`

	if err := os.WriteFile(menuFile, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp menu file failed: %v", err)
	}

	mapping, err := LoadPathMappingFromMenuFile(menuFile)
	if err != nil {
		t.Fatalf("LoadPathMappingFromMenuFile returned error: %v", err)
	}

	if got := mapping["/api/v1/identity/users"]; got != "menu:account_management" {
		t.Fatalf("expected account mapping, got %q", got)
	}

	if got := mapping["/api/v1/identity/users/*"]; got != "menu:account_management" {
		t.Fatalf("expected wildcard account mapping, got %q", got)
	}

	if got := mapping["/api/v1/permission/roles"]; got != "menu:role_permissions" {
		t.Fatalf("expected role mapping, got %q", got)
	}
}

func TestResolveMenuMappingPath(t *testing.T) {
	tempDir := t.TempDir()
	menuFile := filepath.Join(tempDir, "menu.yaml")
	if err := os.WriteFile(menuFile, []byte("menu: []\n"), 0o600); err != nil {
		t.Fatalf("write temp menu file failed: %v", err)
	}

	resolved := resolveMenuMappingPath(menuFile)
	if resolved != menuFile {
		t.Fatalf("expected resolved path %q, got %q", menuFile, resolved)
	}
}

func TestLoadPathMappingFromMenuFile_EmptyPath(t *testing.T) {
	mapping, err := LoadPathMappingFromMenuFile("")
	if err != nil {
		t.Fatalf("expected no error for empty path, got %v", err)
	}

	if len(mapping) != 0 {
		t.Fatalf("expected empty mapping for empty path, got %d", len(mapping))
	}
}
