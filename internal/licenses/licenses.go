package licenses

import _ "embed"

//go:generate go run ../../tools/licensepack

// notices is generated before release and embedded into the final executable.
//
//go:embed generated/THIRD_PARTY_NOTICES.txt
var notices string

func ThirdPartyNotices() string { return notices }
