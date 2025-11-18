package algorithm

import (
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("Basic Push and Pop", func(t *testing.T) {
		stack := NewIntStack()

		// Push
		stack.Push(10)
		stack.Push(20)
		stack.Push(30)

		if stack.Size() != 3 {
			t.Errorf("Expected size 3, got %d", stack.Size())
		}

		// Pop
		if item, ok := stack.Pop(); !ok || item != 30 {
			t.Errorf("Expected 30, got %d", item)
		}

		if item, ok := stack.Pop(); !ok || item != 20 {
			t.Errorf("Expected 20, got %d", item)
		}

		if stack.Size() != 1 {
			t.Errorf("Expected size 1, got %d", stack.Size())
		}
	})

	t.Run("Peek", func(t *testing.T) {
		stack := NewIntStack()

		stack.Push(100)
		stack.Push(200)

		if item, ok := stack.Peek(); !ok || item != 200 {
			t.Errorf("Expected 200, got %d", item)
		}

		// Size should not change after Peek
		if stack.Size() != 2 {
			t.Errorf("Expected size 2 after Peek, got %d", stack.Size())
		}
	})

	t.Run("Empty Stack", func(t *testing.T) {
		stack := NewIntStack()

		if !stack.IsEmpty() {
			t.Error("Expected stack to be empty")
		}

		if _, ok := stack.Pop(); ok {
			t.Error("Expected Pop to fail on empty stack")
		}

		if _, ok := stack.Peek(); ok {
			t.Error("Expected Peek to fail on empty stack")
		}
	})
}

func TestParentheses(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"(())", true},
		{"(()())", true},
		{"(()", false},
		{"())", false},
		{"{[()]}", true},
		{"{[(])}", false},
		{"", true},
	}

	for _, tt := range tests {
		result := isValidParentheses(tt.input)
		if result != tt.expected {
			t.Errorf("isValidParentheses(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestEvaluateRPN(t *testing.T) {
	tests := []struct {
		tokens   []string
		expected int
	}{
		{[]string{"2", "3", "+"}, 5},
		{[]string{"2", "3", "+", "4", "*"}, 20},
		{[]string{"5", "1", "2", "+", "4", "*", "+"}, 17},
	}

	for _, tt := range tests {
		result := evaluateRPN(tt.tokens)
		if result != tt.expected {
			t.Errorf("evaluateRPN(%v) = %d, want %d", tt.tokens, result, tt.expected)
		}
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"a", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		result := reverseString(tt.input)
		if result != tt.expected {
			t.Errorf("reverseString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSimplifyPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/home/user/documents/../pictures", "/home/user/pictures"},
		{"/a/./b/../../c/", "/c"},
		{"/a//b////c/d//././/..", "/a/b/c"},
	}

	for _, tt := range tests {
		result := simplifyPath(tt.input)
		if result != tt.expected {
			t.Errorf("simplifyPath(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
