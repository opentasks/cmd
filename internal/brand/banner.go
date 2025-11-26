// Package brand provides branding assets for the opentask CLI.
package brand

import (
	"fmt"
	"io"
)

// Banner is the ASCII art banner for opentask
const Banner = `
                         __            __  
  ____  ____  ___  ____  / /_____ ____ / /__
 / __ \/ __ \/ _ \/ __ \/ __/ __ '/ __/ //_/
/ /_/ / /_/ /  __/ / / / /_/ /_/ (__  )  <  
\____/ .___/\___/_/ /_/\__/\__,_/____/_/|_| 
    /_/                                     
`

// Version can be set at build time
var Version = "dev"

// PrintBanner writes the banner to the given writer
func PrintBanner(w io.Writer) {
	fmt.Fprint(w, Banner)
}

// PrintBannerWithVersion writes the banner with version info to the given writer
func PrintBannerWithVersion(w io.Writer) {
	fmt.Fprint(w, Banner)
	fmt.Fprintf(w, "  Task management with markdown files (v%s)\n\n", Version)
}
