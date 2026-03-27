package middleware

import (
	"testing"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	permission "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
)

func ptr[T any](v T) *T { return &v }

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "roleIDs", want: "role_ids"},
		{in: "menuID", want: "menu_id"},
		{in: "userProfile", want: "user_profile"},
		{in: "tokenInfo", want: "token_info"},
		{in: "accessToken", want: "access_token"},
		{in: "createdAt", want: "created_at"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := camelToSnake(tt.in)
			if got != tt.want {
				t.Fatalf("camelToSnake(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToSnakeCasePayload_LoginResponseDTO(t *testing.T) {
	perm := permission.PermissionLevel_PERMISSION_LEVEL_FULL
	loginResp := &identity.LoginResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    ptr[int32](0),
			Message: ptr("success"),
		},
		UserProfile: &identity.UserProfileDTO{
			Id:        ptr("u1"),
			Username:  ptr("superadmin"),
			FirstName: ptr("Super"),
			RoleIDs:   []string{"role:superadmin"},
		},
		MenuTree: []*permission.MenuNodeDTO{
			{
				Id:              ptr("system_settings"),
				Name:            ptr("系统设置"),
				HasPermission:   ptr(true),
				PermissionLevel: &perm,
			},
		},
		TokenInfo: &http_base.TokenInfoDTO{
			AccessToken: ptr("token"),
			ExpiresIn:   ptr[int64](1799),
			TokenType:   ptr("Bearer"),
		},
		RoleIDs: []string{"role:superadmin"},
		Permissions: []*permission.MenuPermissionDTO{
			{MenuID: ptr("system_settings"), Permission: &perm},
		},
		Roles: []*identity.RoleInfoDTO{
			{Id: ptr("r1"), Code: ptr("role:superadmin"), Name: ptr("superadmin")},
		},
	}

	payload, err := toSnakeCasePayload(loginResp)
	if err != nil {
		t.Fatalf("toSnakeCasePayload returned error: %v", err)
	}

	root, ok := payload.(map[string]interface{})
	if !ok {
		t.Fatalf("payload type = %T, want map[string]interface{}", payload)
	}

	if _, exists := root["base_resp"]; !exists {
		t.Fatalf("expected key base_resp")
	}
	if _, exists := root["user_profile"]; !exists {
		t.Fatalf("expected key user_profile")
	}
	if _, exists := root["menu_tree"]; !exists {
		t.Fatalf("expected key menu_tree")
	}
	if _, exists := root["token_info"]; !exists {
		t.Fatalf("expected key token_info")
	}
	if _, exists := root["role_ids"]; !exists {
		t.Fatalf("expected key role_ids")
	}
	if _, exists := root["permissions"]; !exists {
		t.Fatalf("expected key permissions")
	}

	if _, exists := root["baseResp"]; exists {
		t.Fatalf("unexpected camelCase key baseResp")
	}
	if _, exists := root["roleIDs"]; exists {
		t.Fatalf("unexpected camelCase key roleIDs")
	}

	userProfile, ok := root["user_profile"].(map[string]interface{})
	if !ok {
		t.Fatalf("user_profile type = %T, want map[string]interface{}", root["user_profile"])
	}
	if _, exists := userProfile["first_name"]; !exists {
		t.Fatalf("expected user_profile.first_name")
	}
	if _, exists := userProfile["role_ids"]; !exists {
		t.Fatalf("expected user_profile.role_ids")
	}

	tokenInfo, ok := root["token_info"].(map[string]interface{})
	if !ok {
		t.Fatalf("token_info type = %T, want map[string]interface{}", root["token_info"])
	}
	if _, exists := tokenInfo["access_token"]; !exists {
		t.Fatalf("expected token_info.access_token")
	}
	if _, exists := tokenInfo["expires_in"]; !exists {
		t.Fatalf("expected token_info.expires_in")
	}

	permissionsArr, ok := root["permissions"].([]interface{})
	if !ok || len(permissionsArr) == 0 {
		t.Fatalf("permissions type = %T, want non-empty []interface{}", root["permissions"])
	}
	permission0, ok := permissionsArr[0].(map[string]interface{})
	if !ok {
		t.Fatalf("permissions[0] type = %T, want map[string]interface{}", permissionsArr[0])
	}
	if _, exists := permission0["menu_id"]; !exists {
		t.Fatalf("expected permissions[0].menu_id")
	}
	if _, exists := permission0["menuID"]; exists {
		t.Fatalf("unexpected permissions[0].menuID")
	}

	menuTree, ok := root["menu_tree"].([]interface{})
	if !ok || len(menuTree) == 0 {
		t.Fatalf("menu_tree type = %T, want non-empty []interface{}", root["menu_tree"])
	}
	menuNode0, ok := menuTree[0].(map[string]interface{})
	if !ok {
		t.Fatalf("menu_tree[0] type = %T, want map[string]interface{}", menuTree[0])
	}
	if _, exists := menuNode0["has_permission"]; !exists {
		t.Fatalf("expected menu_tree[0].has_permission")
	}
	if _, exists := menuNode0["permission_level"]; !exists {
		t.Fatalf("expected menu_tree[0].permission_level")
	}
}
