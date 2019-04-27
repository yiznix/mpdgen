// Copyright 2019 yiznix.com. All rights reserved.
// Use of this source code is governed by the license that can be found in LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/yiznix/mpdgen"
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
	mpdgen.MPDGen(*srcFile, *destDir, *segment)
}
