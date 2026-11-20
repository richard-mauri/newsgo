package dashboard

// One semantic version feeds the banner, CLI and cache-busting URLs.
// Increment only for an approved release, not each development build.
const appVersion = "0.2.58"

func Version() string { return appVersion }

// BuildNumber is injected by devbuild.sh and never committed as a local count.
var BuildNumber = "000000"
func BuildLabel() string { return "Build " + BuildNumber }
