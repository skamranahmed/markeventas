package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testData struct {
	input          string
	expectedOutput string
	expectedError  error // default value is nil
}

func TestProcessMonthName(t *testing.T) {
	tests := []testData{
		{input: "Jan", expectedOutput: January},
		{input: "jan", expectedOutput: January},
		{input: "January", expectedOutput: January},
		{input: "JAN", expectedOutput: January},
		{input: "JANUARY", expectedOutput: January},
		{input: "Janauaray", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "Feb", expectedOutput: February},
		{input: "feb", expectedOutput: February},
		{input: "February", expectedOutput: February},
		{input: "FEB", expectedOutput: February},
		{input: "FEBRUARY", expectedOutput: February},
		{input: "Fefhfhnvj", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "Mar", expectedOutput: March},
		{input: "mar", expectedOutput: March},
		{input: "March", expectedOutput: March},
		{input: "MAR", expectedOutput: March},
		{input: "MARCH", expectedOutput: March},
		{input: "Marc", expectedOutput: "", expectedError: ErrInvalidMonth},
		{input: "Maajkfh", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "Apr", expectedOutput: April},
		{input: "apr", expectedOutput: April},
		{input: "April", expectedOutput: April},
		{input: "APR", expectedOutput: April},
		{input: "APRIL", expectedOutput: April},
		{input: "Apri", expectedOutput: "", expectedError: ErrInvalidMonth},
		{input: "Apsflfjfu", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "May", expectedOutput: May},
		{input: "may", expectedOutput: May},
		{input: "MAY", expectedOutput: May},
		{input: "m", expectedOutput: "", expectedError: ErrInvalidMonth},
		{input: "Ma", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "Jun", expectedOutput: June},
		{input: "jun", expectedOutput: June},
		{input: "June", expectedOutput: June},
		{input: "JUNE", expectedOutput: June},
		{input: "ju", expectedOutput: "", expectedError: ErrInvalidMonth},
		{input: "jskhksi", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "Jul", expectedOutput: July},
		{input: "jul", expectedOutput: July},
		{input: "July", expectedOutput: July},
		{input: "JULY", expectedOutput: July},

		{input: "Aug", expectedOutput: August},
		{input: "aug", expectedOutput: August},
		{input: "augst", expectedOutput: August},
		{input: "agst", expectedOutput: August},
		{input: "August", expectedOutput: August},
		{input: "AUGUST", expectedOutput: August},

		{input: "Sep", expectedOutput: September},
		{input: "sep", expectedOutput: September},
		{input: "sept", expectedOutput: September},
		{input: "SEPT", expectedOutput: September},
		{input: "september", expectedOutput: September},
		{input: "spfpsgo", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "oct", expectedOutput: October},
		{input: "OCT", expectedOutput: October},
		{input: "october", expectedOutput: October},
		{input: "OCTOBER", expectedOutput: October},
		{input: "octo", expectedOutput: "", expectedError: ErrInvalidMonth},
		{input: "ocodkog", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "nov", expectedOutput: November},
		{input: "November", expectedOutput: November},
		{input: "NOV", expectedOutput: November},
		{input: "novo", expectedOutput: "", expectedError: ErrInvalidMonth},

		{input: "December", expectedOutput: December},
		{input: "dec", expectedOutput: December},
		{input: "DEC", expectedOutput: December},
		{input: "DECEMBER", expectedOutput: December},
		{input: "decem", expectedOutput: "", expectedError: ErrInvalidMonth},
	}

	for _, test := range tests {
		got, err := processMonthName(test.input)
		msg := fmt.Sprintf("input: %s, got: %s, expected: %s, err: %v", test.input, got, test.expectedOutput, err)
		require.Equal(t, test.expectedError, err, msg)
		require.Equal(t, test.expectedOutput, got, msg)
	}
}
