package idam

var passwordValidationCases = []struct {
	name      string
	password  string
	wantValid bool
}{
	{"too short (11 chars)", "abcdefghijk", false},
	{"minimum length (12 chars, no complexity)", "abcdefghijkl", true},
	{"long passphrase no complexity", "correct horse battery staple", true},
	{"too long (65 chars)", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
	{"non-ASCII character", "abcdefghijkl\xe9", false},
	{"empty", "", false},
}
