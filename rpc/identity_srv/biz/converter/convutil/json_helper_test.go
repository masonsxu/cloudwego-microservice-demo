package convutil

import (
	"testing"
)

// ==================== StringSliceToJSON ====================

func TestStringSliceToJSON_Empty(t *testing.T) {
	result := StringSliceToJSON([]string{})
	if result != "" {
		t.Errorf("StringSliceToJSON([]) = %q, want empty string", result)
	}
}

func TestStringSliceToJSON_Nil(t *testing.T) {
	result := StringSliceToJSON(nil)
	if result != "" {
		t.Errorf("StringSliceToJSON(nil) = %q, want empty string", result)
	}
}

func TestStringSliceToJSON_SingleElement(t *testing.T) {
	result := StringSliceToJSON([]string{"hello"})
	expected := `["hello"]`
	if result != expected {
		t.Errorf("StringSliceToJSON([hello]) = %q, want %q", result, expected)
	}
}

func TestStringSliceToJSON_MultipleElements(t *testing.T) {
	result := StringSliceToJSON([]string{"foo", "bar", "baz"})
	expected := `["foo","bar","baz"]`
	if result != expected {
		t.Errorf("StringSliceToJSON() = %q, want %q", result, expected)
	}
}

func TestStringSliceToJSON_SpecialCharacters(t *testing.T) {
	result := StringSliceToJSON([]string{"hello\"world", "foo\nbar"})
	expected := `["hello\"world","foo\nbar"]`
	if result != expected {
		t.Errorf("StringSliceToJSON() = %q, want %q", result, expected)
	}
}

func TestStringSliceToJSON_Unicode(t *testing.T) {
	result := StringSliceToJSON([]string{"你好", "世界"})
	expected := `["你好","世界"]`
	if result != expected {
		t.Errorf("StringSliceToJSON() = %q, want %q", result, expected)
	}
}

// ==================== JSONToStringSlice ====================

func TestJSONToStringSlice_Empty(t *testing.T) {
	result := JSONToStringSlice("")
	if result != nil {
		t.Errorf("JSONToStringSlice(\"\") = %v, want nil", result)
	}
}

func TestJSONToStringSlice_EmptyArray(t *testing.T) {
	result := JSONToStringSlice("[]")
	if len(result) != 0 {
		t.Errorf("JSONToStringSlice(\"[]\") = %v, want empty slice", result)
	}
}

func TestJSONToStringSlice_SingleElement(t *testing.T) {
	result := JSONToStringSlice(`["hello"]`)
	if len(result) != 1 || result[0] != "hello" {
		t.Errorf("JSONToStringSlice() = %v, want [hello]", result)
	}
}

func TestJSONToStringSlice_MultipleElements(t *testing.T) {
	result := JSONToStringSlice(`["foo","bar","baz"]`)
	expected := []string{"foo", "bar", "baz"}
	if !sliceEqual(result, expected) {
		t.Errorf("JSONToStringSlice() = %v, want %v", result, expected)
	}
}

func TestJSONToStringSlice_InvalidJSON(t *testing.T) {
	result := JSONToStringSlice("not json")
	if result != nil {
		t.Errorf("JSONToStringSlice(invalid) = %v, want nil", result)
	}
}

func TestJSONToStringSlice_WrongType(t *testing.T) {
	result := JSONToStringSlice(`{"key": "value"}`)
	if result != nil {
		t.Errorf("JSONToStringSlice(object) = %v, want nil", result)
	}
}

// ==================== ULIDSliceToJSON / JSONToULIDSlice ====================

func TestULIDSliceToJSON_Empty(t *testing.T) {
	result := ULIDSliceToJSON([]string{})
	if result != "" {
		t.Errorf("ULIDSliceToJSON([]) = %q, want empty string", result)
	}
}

func TestULIDSliceToJSON_ValidULIDs(t *testing.T) {
	ulids := []string{"01ARZ3NDEKTSV4RRFFQ69G5FAV", "01ARZ3NDEKTSV4RRFFQ69G5FAW"}
	result := ULIDSliceToJSON(ulids)
	expected := `["01ARZ3NDEKTSV4RRFFQ69G5FAV","01ARZ3NDEKTSV4RRFFQ69G5FAW"]`
	if result != expected {
		t.Errorf("ULIDSliceToJSON() = %q, want %q", result, expected)
	}
}

func TestJSONToULIDSlice_Empty(t *testing.T) {
	result := JSONToULIDSlice("")
	if result != nil {
		t.Errorf("JSONToULIDSlice(\"\") = %v, want nil", result)
	}
}

func TestJSONToULIDSlice_Valid(t *testing.T) {
	result := JSONToULIDSlice(`["01ARZ3NDEKTSV4RRFFQ69G5FAV"]`)
	if len(result) != 1 || result[0] != "01ARZ3NDEKTSV4RRFFQ69G5FAV" {
		t.Errorf("JSONToULIDSlice() = %v, want [01ARZ3NDEKTSV4RRFFQ69G5FAV]", result)
	}
}

// ==================== StringPtr / StringValue ====================

func TestStringPtr_Empty(t *testing.T) {
	result := StringPtr("")
	if result != nil {
		t.Errorf("StringPtr(\"\") = %v, want nil", result)
	}
}

func TestStringPtr_NonEmpty(t *testing.T) {
	result := StringPtr("hello")
	if result == nil || *result != "hello" {
		t.Errorf("StringPtr(\"hello\") = %v, want *\"hello\"", result)
	}
}

func TestStringValue_Nil(t *testing.T) {
	result := StringValue(nil)
	if result != "" {
		t.Errorf("StringValue(nil) = %q, want empty string", result)
	}
}

func TestStringValue_NonNil(t *testing.T) {
	s := "hello"
	result := StringValue(&s)
	if result != "hello" {
		t.Errorf("StringValue(&\"hello\") = %q, want \"hello\"", result)
	}
}

// ==================== Int32Ptr / Int32Value ====================

func TestInt32Ptr_Zero(t *testing.T) {
	result := Int32Ptr(0)
	if result != nil {
		t.Errorf("Int32Ptr(0) = %v, want nil", result)
	}
}

func TestInt32Ptr_NonZero(t *testing.T) {
	result := Int32Ptr(42)
	if result == nil || *result != 42 {
		t.Errorf("Int32Ptr(42) = %v, want *42", result)
	}
}

func TestInt32Ptr_Negative(t *testing.T) {
	result := Int32Ptr(-1)
	if result == nil || *result != -1 {
		t.Errorf("Int32Ptr(-1) = %v, want *-1", result)
	}
}

func TestInt32Value_Nil(t *testing.T) {
	result := Int32Value(nil)
	if result != 0 {
		t.Errorf("Int32Value(nil) = %d, want 0", result)
	}
}

func TestInt32Value_NonNil(t *testing.T) {
	var i int32 = 42
	result := Int32Value(&i)
	if result != 42 {
		t.Errorf("Int32Value(&42) = %d, want 42", result)
	}
}

// ==================== Int64Ptr / Int64Value ====================

func TestInt64Ptr_Zero(t *testing.T) {
	result := Int64Ptr(0)
	if result != nil {
		t.Errorf("Int64Ptr(0) = %v, want nil", result)
	}
}

func TestInt64Ptr_NonZero(t *testing.T) {
	result := Int64Ptr(1703472000000)
	if result == nil || *result != 1703472000000 {
		t.Errorf("Int64Ptr(timestamp) = %v, want *1703472000000", result)
	}
}

func TestInt64Value_Nil(t *testing.T) {
	result := Int64Value(nil)
	if result != 0 {
		t.Errorf("Int64Value(nil) = %d, want 0", result)
	}
}

func TestInt64Value_NonNil(t *testing.T) {
	var i int64 = 1703472000000
	result := Int64Value(&i)
	if result != 1703472000000 {
		t.Errorf("Int64Value(&timestamp) = %d, want 1703472000000", result)
	}
}

// ==================== BoolPtr / BoolValue ====================

func TestBoolPtr_True(t *testing.T) {
	result := BoolPtr(true)
	if result == nil || *result != true {
		t.Errorf("BoolPtr(true) = %v, want *true", result)
	}
}

func TestBoolPtr_False(t *testing.T) {
	result := BoolPtr(false)
	if result == nil || *result != false {
		t.Errorf("BoolPtr(false) = %v, want *false", result)
	}
}

func TestBoolValue_Nil(t *testing.T) {
	result := BoolValue(nil)
	if result != false {
		t.Errorf("BoolValue(nil) = %v, want false", result)
	}
}

func TestBoolValue_True(t *testing.T) {
	b := true
	result := BoolValue(&b)
	if result != true {
		t.Errorf("BoolValue(&true) = %v, want true", result)
	}
}

func TestBoolValue_False(t *testing.T) {
	b := false
	result := BoolValue(&b)
	if result != false {
		t.Errorf("BoolValue(&false) = %v, want false", result)
	}
}

// ==================== TrimSpace ====================

func TestTrimSpace_Empty(t *testing.T) {
	result := TrimSpace("")
	if result != "" {
		t.Errorf("TrimSpace(\"\") = %q, want empty string", result)
	}
}

func TestTrimSpace_OnlySpaces(t *testing.T) {
	result := TrimSpace("   ")
	if result != "" {
		t.Errorf("TrimSpace(\"   \") = %q, want empty string", result)
	}
}

func TestTrimSpace_LeadingTrailing(t *testing.T) {
	result := TrimSpace("  hello  ")
	if result != "hello" {
		t.Errorf("TrimSpace(\"  hello  \") = %q, want \"hello\"", result)
	}
}

func TestTrimSpace_Tabs(t *testing.T) {
	result := TrimSpace("\t\nhello\t\n")
	if result != "hello" {
		t.Errorf("TrimSpace(with tabs) = %q, want \"hello\"", result)
	}
}

// ==================== IsEmptyString ====================

func TestIsEmptyString_Empty(t *testing.T) {
	if !IsEmptyString("") {
		t.Error("IsEmptyString(\"\") = false, want true")
	}
}

func TestIsEmptyString_OnlySpaces(t *testing.T) {
	if !IsEmptyString("   ") {
		t.Error("IsEmptyString(\"   \") = false, want true")
	}
}

func TestIsEmptyString_OnlyTabs(t *testing.T) {
	if !IsEmptyString("\t\n") {
		t.Error("IsEmptyString(tabs) = false, want true")
	}
}

func TestIsEmptyString_NonEmpty(t *testing.T) {
	if IsEmptyString("hello") {
		t.Error("IsEmptyString(\"hello\") = true, want false")
	}
}

func TestIsEmptyString_SpacesAround(t *testing.T) {
	if IsEmptyString("  hello  ") {
		t.Error("IsEmptyString(\"  hello  \") = true, want false")
	}
}

// ==================== 往返测试 (Roundtrip) ====================

func TestStringSlice_Roundtrip(t *testing.T) {
	original := []string{"foo", "bar", "baz"}
	json := StringSliceToJSON(original)
	result := JSONToStringSlice(json)

	if !sliceEqual(original, result) {
		t.Errorf("Roundtrip failed: original=%v, result=%v", original, result)
	}
}

func TestULIDSlice_Roundtrip(t *testing.T) {
	original := []string{"01ARZ3NDEKTSV4RRFFQ69G5FAV", "01ARZ3NDEKTSV4RRFFQ69G5FAW"}
	json := ULIDSliceToJSON(original)
	result := JSONToULIDSlice(json)

	if !sliceEqual(original, result) {
		t.Errorf("Roundtrip failed: original=%v, result=%v", original, result)
	}
}

// ==================== 辅助函数 ====================

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
