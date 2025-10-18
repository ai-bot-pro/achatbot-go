package functions

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAdapteOpenAIToolSchema(t *testing.T) {
	tests := []struct {
		name          string // description of this test case
		schemas       []map[string]any
		exceptSchemas []map[string]any
		wantErr       bool
	}{
		{
			name:          "empty schemas",
			schemas:       []map[string]any{},
			exceptSchemas: []map[string]any{},
			wantErr:       false,
		},
		{
			name: "single search tool schema",
			schemas: []map[string]any{
				SearchToolSchema,
			},
			exceptSchemas: []map[string]any{
				SearchToolSchema,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gots, gotsErr := AdapteOpenAIToolSchema(tt.schemas)
			if gotsErr != nil {
				if !tt.wantErr {
					t.Errorf("AdapteOpenAIToolSchema() failed: %v", gotsErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AdapteOpenAIToolSchema() succeeded unexpectedly")
			}
			// Basic length check
			if len(gots) != len(tt.schemas) {
				t.Errorf("AdapteOpenAIToolSchema() length = %v, want %v", len(gots), len(tt.schemas))
			}
			if len(gots) > 0 {
				for i, got := range gots {
					var gotMap map[string]any
					bytes, _ := got.MarshalJSON()
					json.Unmarshal(bytes, &gotMap)

					// 为确保类型一致性，将期望值也进行 JSON 序列化和反序列化
					expectedBytes, _ := json.Marshal(tt.exceptSchemas[i])
					var expectedMap map[string]any
					json.Unmarshal(expectedBytes, &expectedMap)

					ok := reflect.DeepEqual(gotMap, expectedMap)
					if !ok {
						t.Errorf("%s mismatch: gots:\n%+v \nwant:\n%+v", tt.name, gotMap, expectedMap)
					}
				}
			}
		})
	}
}

func TestAdapteOllamaToolSchema(t *testing.T) {
	tests := []struct {
		name          string // description of this test case
		schemas       []map[string]any
		exceptSchemas []map[string]any
		wantErr       bool
	}{
		{
			name:          "empty schemas",
			schemas:       []map[string]any{},
			exceptSchemas: []map[string]any{},
			wantErr:       false,
		},
		{
			name: "single search tool schema",
			schemas: []map[string]any{
				OllamaAPISearchToolSchema,
			},
			exceptSchemas: []map[string]any{
				SearchToolSchema,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gots, gotsErr := AdapteOllamaToolSchema(tt.schemas)
			if gotsErr != nil {
				if !tt.wantErr {
					t.Errorf("AdapteOllamaToolSchema() failed: %v", gotsErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AdapteOllamaToolSchema() succeeded unexpectedly")
			}
			// Basic length check
			if len(gots) != len(tt.schemas) {
				t.Errorf("AdapteOllamaToolSchema() length = %v, want %v", len(gots), len(tt.schemas))
			}
			if len(gots) > 0 {
				for i, got := range gots {
					gotBytes, _ := json.Marshal(got)
					var gotMap map[string]any
					json.Unmarshal(gotBytes, &gotMap)

					// 为确保类型一致性，将期望值也进行 JSON 序列化和反序列化
					expectedBytes, _ := json.Marshal(tt.exceptSchemas[i])
					var expectedMap map[string]any
					json.Unmarshal(expectedBytes, &expectedMap)

					ok := reflect.DeepEqual(gotMap, expectedMap)
					if !ok {
						t.Errorf("%s mismatch: gots:\n%+v \nwant:\n%+v", tt.name, gotMap, expectedMap)
					}
				}
			}
		})
	}
}
