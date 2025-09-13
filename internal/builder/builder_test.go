package builder

import (
	"testing"
)

func TestBuilder_NewBuilder(t *testing.T) {
	b := NewBuilder()

	if b == nil {
		t.Fatal("Expected NewBuilder() to return non-nil")
	}
}

func TestBuilder_AddInteger(t *testing.T) {
	b := NewBuilder()

	// Test that AddInteger returns the builder for chaining
	result := b.AddInteger(42)
	if result != b {
		t.Error("Expected AddInteger() to return the same builder instance")
	}

	// Test multiple additions for chaining
	b2 := b.
		AddInteger(1).
		AddInteger(17).
		AddInteger(42)

	if b2 != b {
		t.Error("Expected chaining to work correctly")
	}
}

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Builder)
		wantErr bool
	}{
		{
			name: "empty builder",
			setup: func(b *Builder) {
				// No additions
			},
			wantErr: false,
		},
		{
			name: "single integer",
			setup: func(b *Builder) {
				b.AddInteger(42)
			},
			wantErr: false,
		},
		{
			name: "multiple integers",
			setup: func(b *Builder) {
				b.AddInteger(1).AddInteger(17).AddInteger(42)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			tt.setup(b)

			bs, err := b.Build()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected Build() to return an error")
				}
				if bs != nil {
					t.Error("Expected Build() to return nil bitstring on error")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if bs == nil {
				t.Fatal("Expected Build() to return non-nil bitstring")
			}

			// Basic validation of created bitstring
			if bs.Length() == 0 && tt.name != "empty builder" {
				t.Errorf("Expected non-zero length for %s", tt.name)
			}

			// Should be binary since we're adding integers (default 8 bits each)
			if tt.name != "empty builder" && !bs.IsBinary() {
				t.Error("Expected bitstring to be binary")
			}
		})
	}
}

func TestBuilder_BuildContent(t *testing.T) {
	// Test specific content generation
	b := NewBuilder().
		AddInteger(1).
		AddInteger(17).
		AddInteger(42)

	bs, err := b.Build()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bs.Length() != 24 { // 3 integers * 8 bits = 24 bits
		t.Errorf("Expected bitstring length 24, got %d", bs.Length())
	}

	if !bs.IsBinary() {
		t.Error("Expected bitstring to be binary")
	}

	bytes := bs.ToBytes()
	if len(bytes) != 3 {
		t.Fatalf("Expected 3 bytes, got %d", len(bytes))
	}

	if bytes[0] != 1 || bytes[1] != 17 || bytes[2] != 42 {
		t.Errorf("Expected [1, 17, 42], got %v", bytes)
	}
}

func TestBuilder_EmptyBuild(t *testing.T) {
	b := NewBuilder()
	bs, err := b.Build()

	if err != nil {
		t.Errorf("Expected no error for empty build, got %v", err)
	}

	if bs == nil {
		t.Fatal("Expected non-nil bitstring for empty build")
	}

	if bs.Length() != 0 {
		t.Errorf("Expected empty bitstring length 0, got %d", bs.Length())
	}

	if !bs.IsEmpty() {
		t.Error("Expected empty bitstring to be empty")
	}

	if !bs.IsBinary() {
		t.Error("Expected empty bitstring to be binary")
	}

	bytes := bs.ToBytes()
	if len(bytes) != 0 {
		t.Errorf("Expected empty bytes, got %v", bytes)
	}
}
