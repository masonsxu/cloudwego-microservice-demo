package casbin_middleware

import "testing"

func TestLoadConfigFromEnv_SuperAdminConfig(t *testing.T) {
	t.Setenv("CASBIN_SUPERADMIN_BYPASS_ENABLED", "false")
	t.Setenv("CASBIN_SUPERADMIN_SUBJECTS", "role:superadmin, role:platform_admin")

	config := LoadConfigFromEnv()
	if config.SuperAdminBypassEnabled {
		t.Error("SuperAdminBypassEnabled should be false when env is false")
	}

	if len(config.SuperAdminSubjects) != 2 {
		t.Fatalf("expected 2 super admin subjects, got %d", len(config.SuperAdminSubjects))
	}

	if config.SuperAdminSubjects[0] != "role:superadmin" {
		t.Errorf("unexpected first subject: %s", config.SuperAdminSubjects[0])
	}

	if config.SuperAdminSubjects[1] != "role:platform_admin" {
		t.Errorf("unexpected second subject: %s", config.SuperAdminSubjects[1])
	}
}
