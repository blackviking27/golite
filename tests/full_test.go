package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blackviking27/golite/operations"
	"github.com/blackviking27/golite/types"
)

func TestRowSerialization(t *testing.T) {
	// Define our test cases
	tests := []struct {
		name     string
		input    types.Row
		expected types.Row // Expected might differ if we truncate long strings
	}{
		{
			name:     "Normal User",
			input:    types.Row{ID: 1, Username: "janedoe", Email: "jane@example.com"},
			expected: types.Row{ID: 1, Username: "janedoe", Email: "jane@example.com"},
		},
		{
			name:     "Empty Fields",
			input:    types.Row{ID: 2, Username: "", Email: ""},
			expected: types.Row{ID: 2, Username: "", Email: ""},
		},
		{
			name: "Max Length Username (32 chars)",
			// Create a string exactly 32 chars long
			input:    types.Row{ID: 3, Username: strings.Repeat("a", 32), Email: "test@test.com"},
			expected: types.Row{ID: 3, Username: strings.Repeat("a", 32), Email: "test@test.com"},
		},
		{
			name: "Truncate Long Username",
			// Input is 33 chars, expectation is 32 chars
			input:    types.Row{ID: 4, Username: strings.Repeat("a", 33), Email: "test@test.com"},
			expected: types.Row{ID: 4, Username: strings.Repeat("a", 32), Email: "test@test.com"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup: Create a clean buffer (simulating a row slot in a page)
			buffer := make([]byte, types.RowSize)

			// 2. Action: Serialize
			operations.SerializeRow(&tc.input, buffer)

			// 3. Action: Deserialize
			var output types.Row
			operations.DeserializeRow(buffer, &output)

			// 4. Assertion: Check ID
			if output.ID != tc.expected.ID {
				t.Errorf("ID mismatch: got %d, want %d", output.ID, tc.expected.ID)
			}
			// 5. Assertion: Check Username
			if output.Username != tc.expected.Username {
				t.Errorf("Username mismatch:\nGot:  %q\nWant: %q", output.Username, tc.expected.Username)
			}
			// 6. Assertion: Check Email
			if output.Email != tc.expected.Email {
				t.Errorf("Email mismatch:\nGot:  %q\nWant: %q", output.Email, tc.expected.Email)
			}
		})
	}
}

func TestInsertStatement(t *testing.T) {
	// Prepare statment should fail for negative id for insert statement
	command := "insert -1 johndoe johndoe@gmail.com"

	var statement1 types.Statement
	status := operations.PrepareStatement(command, &statement1)

	if status != types.PREPARE_NEGATIVE_ID_FAILURE {
		t.Error("Failed to validate negative id in insert statement")
	}

	// Testing the max character length check
	command = "insert 1 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa test@test.com"
	var statement2 types.Statement
	status = operations.PrepareStatement(command, &statement2)

	if status != types.PREPARE_STRING_TOO_LONG_FAILURE {
		t.Error("Failed to validate max length in statement")
	}

}

// 2. Test Max Rows Calculation
// This ensures our math for Page Size and Limits is correct.
func TestTableConstants(t *testing.T) {
	// 4096 / 291 = 14.07... -> Should be 14 rows per page
	expectedRowsPerPage := uint32(14)
	if types.RowsPerPage != expectedRowsPerPage {
		t.Errorf("RowsPerPage calculation wrong. Got %d, want %d", types.RowsPerPage, expectedRowsPerPage)
	}

	// 14 * 100 = 1400 max rows
	expectedMaxRows := expectedRowsPerPage * types.TableMaxPages
	if types.TableMaxRows != expectedMaxRows {
		t.Errorf("TableMaxRows calculation wrong. Got %d, want %d", types.TableMaxRows, expectedMaxRows)
	}
}

// 3. Test RowSlot Page Allocation
// This simulates inserting rows to verify pages are allocated when needed.
func TestRowSlotAllocation(t *testing.T) {
	table := operations.NewTable()

	// Case 1: Fetch Row 0. Should allocate Page 0.
	_ = table.RowSlot(0)
	if table.Pages[0] == nil {
		t.Error("RowSlot(0) failed to allocate Page 0")
	}
	if len(table.Pages[0]) != types.PageSize {
		t.Errorf("Page 0 has wrong size. Got %d, want %d", len(table.Pages[0]), types.PageSize)
	}

	// Case 2: Fetch Row 13 (Last row of Page 0). Should NOT allocate Page 1.
	_ = table.RowSlot(13)
	if table.Pages[1] != nil {
		t.Error("RowSlot(13) allocated Page 1 prematurely!")
	}

	// Case 3: Fetch Row 14 (First row of Page 1). Should allocate Page 1.
	_ = table.RowSlot(14)
	if table.Pages[1] == nil {
		t.Error("RowSlot(14) failed to allocate Page 1")
	}
}

func TestTableCapacity(t *testing.T) {
	// 1. Create a fresh table
	table := operations.NewTable()

	// Calculate max rows explicitly for the test (1400 rows)
	maxRows := types.TableMaxRows

	t.Logf("Testing storage of %d rows...", maxRows)

	// 2. Fill the database to its absolute limit
	for i := 0; i < maxRows; i++ {
		// Create a unique row based on the index
		row := types.Row{
			ID:       uint32(i),
			Username: fmt.Sprintf("User%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
		}

		// Get the slot.
		// If your page allocation logic is wrong, this might panic near i=1386 (start of last page)
		dest := table.RowSlot(uint32(i))

		operations.SerializeRow(&row, dest)
	}

	// 3. Verify the Last Row (Boundary Check)
	// We read back the very last possible row to ensure the last page was written correctly.
	lastIndex := maxRows - 1
	src := table.RowSlot(uint32(lastIndex))

	var lastRow types.Row
	operations.DeserializeRow(src, &lastRow)

	// Check if the ID matches what we put in
	if lastRow.ID != uint32(lastIndex) {
		t.Errorf("Failed to fetch the last row (ID %d). Got ID %d", lastIndex, lastRow.ID)
	}

	// 4. (Optional) Verify Out of Bounds Panic
	// This ensures we can't write row 1401
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when writing past the limit!")
		}
	}()

	// Try to access the first invalid row index. This SHOULD panic.
	_ = table.RowSlot(uint32(maxRows))
}
