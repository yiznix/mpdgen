# mpdgen
Utilities to generate MPEG DASH manifest files and segment files.

This pacakge contains utilities which call ffmpeg and MP4Box commands on Linux to generate DASH manifestfiles and audio/video segment files. The generated files can be used for streaming videos on demand.

The utitlities in this package are mainly for testing and experimenting purpose. Use with cautions in production environment.

## Prerequisites

The utilies run on Linux OS and have been tested under Ubuntu16.04 and Ubuntu18.04. ffmpeg and MP4Box must be preinstalled. For ffmpeg please refer to its [installation guide](https://trac.ffmpeg.org/wiki/CompilationGuide/Ubuntu). For MP4Box, install with "apt-get install gpac" command.

## Usage

Generate DASH manifest and setments using MPDGen().

```go
package main

import (
    "flag"
    "log"
)

var (
    srcFile = flag.String("file", "", "The path to the soruce video.")
    destDir = flag.String("dest", "", "The directory where the manifest and segments are stored.")
    segment = flag.String("segment", "1000", "Dash segment value")
)

func main() {
    flag.Parse()

    if *srcFile == "" || *destDir == "" {
        log.Fatalf("The --file and --dest flags must be provided.")
    }
    MPDGen(*srcFile, *destDir, *segment)
}
```

To play the video in browser, you can use the [dash.js libaray](https://github.com/Dash-Industry-Forum/dash.js).
