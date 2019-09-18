package version

var (
	version string
	arch    string
	os      string
)

func Version() string {
	return version
}

func Os() string {
	return os
}

func Arch() string {
	return arch
}
