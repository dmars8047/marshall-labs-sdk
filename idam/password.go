package idam

import "github.com/dmars8047/strval"

func validatePassword(value, fieldName string) strval.StringValidationResult {
	return strval.ValidateStringWithName(value, fieldName,
		strval.MustNotBeEmpty(),
		strval.MustHaveMinLengthOf(MinPasswordLength),
		strval.MustHaveMaxLengthOf(MaxPasswordLength),
		strval.MustOnlyContainPrintableCharacters(),
		strval.MustOnlyContainASCIICharacters())
}
