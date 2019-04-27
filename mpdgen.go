// Copyright 2019 yiznix.com. All rights reserved.
// Use of this source code is governed by the license that can be found in LICENSE file.

package mpdgen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const (
	// FFMPEG is the location of ffmpeg binary on Linux/Ubuntu.
	FFMPEG = "/usr/local/bin/ffmpeg"

	// MP4BOX is the location of MP4Box binary on Linux/Ubuntu.
	MP4BOX = "/usr/bin/MP4Box"

	mpdTpl = `<?xml version="1.0" encoding="UTF-8" ?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" minBufferTime="%s" type="static" mediaPresentationDuration="%s" profiles="urn:mpeg:dash:profile:isoff-live:2011">
  <Period> 
    <AdaptationSet segmentAlignment="true" mimeType="video/mp4" codecs="%s" par="%s">
      <Representation width="%s" height="%s" frameRate="%s" id="1" bandwidth="%s">
        <SegmentTemplate timescale="%s" duration="%s" startNumber="1" media="video$Number$.m4s" initialization="videoinit.mp4"/>
      </Representation>
    </AdaptationSet>
    <AdaptationSet mimeType="audio/mp4" codecs="%s" streamName="stereo">
      <Representation audioSamplingRate="%s" id="2" bandwidth="%s">
        <SegmentTemplate timescale="%s"  startNumber="1" media="audio$Number$.m4s" initialization="audioinit.mp4" duration="%s"/>
      </Representation>
    </AdaptationSet>
  </Period>
</MPD>`
)

type videoAttr struct {
	minBufferTime             string
	mediaPresentationDuration string
	codecs                    string
	par                       string
	width                     string
	height                    string
	frameRate                 string
	bandwidth                 string
	timescale                 string
	duration                  string
}

type audioAttr struct {
	minBufferTime             string
	mediaPresentationDuration string
	codecs                    string
	audioSamplingRate         string
	bandwidth                 string
	timescale                 string
	duration                  string
}

func getVideoAttr(videoMPD string) (*videoAttr, error) {
	// minBufferTime
	minBufferTimeRe := regexp.MustCompile(`minBufferTime="[PYMDTHMS0-9.]*"`)
	minBufferTimeSlice := minBufferTimeRe.FindAllString(videoMPD, -1)
	if len(minBufferTimeSlice) != 1 {
		return nil, fmt.Errorf("minBufferTimeSlice not found")
	}
	minBufferTime := minBufferTimeSlice[0]
	minBufferTime = minBufferTime[len(`minBufferTime="`) : len(minBufferTime)-1]

	// mediaPresentationDuration
	mediaPresentationDurationRe := regexp.MustCompile(`mediaPresentationDuration="[PYMDTHMS0-9.]*"`)
	mediaPresentationDurationSlice := mediaPresentationDurationRe.FindAllString(videoMPD, -1)
	if len(mediaPresentationDurationSlice) != 1 {
		return nil, fmt.Errorf("mediaPresentationDuration not found")
	}
	mediaPresentationDuration := mediaPresentationDurationSlice[0]
	mediaPresentationDuration = mediaPresentationDuration[len(`mediaPresentationDuration="`) : len(mediaPresentationDuration)-1]

	// codecs
	codecsRe := regexp.MustCompile(`codecs="[A-Za-z0-9._-]*"`)
	codecsSlice := codecsRe.FindAllString(videoMPD, -1)
	if len(codecsSlice) != 1 {
		return nil, fmt.Errorf("codecs not found")
	}
	codecs := codecsSlice[0]
	codecs = codecs[len(`codecs="`) : len(codecs)-1]

	// par
	par := ""
	parRe := regexp.MustCompile(`par="[0-9:]*"`)
	parSlice := parRe.FindAllString(videoMPD, -1)
	if len(parSlice) != 1 {
		return nil, fmt.Errorf("par not found")
	}
	par = parSlice[0]
	par = par[len(`par="`) : len(par)-1]

	// width
	widthRe := regexp.MustCompile(` width="[0-9]*"`)
	widthSlice := widthRe.FindAllString(videoMPD, -1)
	if len(widthSlice) != 1 {
		return nil, fmt.Errorf("width not found")
	}
	width := widthSlice[0]
	width = width[len(` width="`) : len(width)-1]

	// height
	heightRe := regexp.MustCompile(`height="[0-9]*"`)
	heightSlice := heightRe.FindAllString(videoMPD, -1)
	if len(heightSlice) != 1 {
		return nil, fmt.Errorf("height not found")
	}
	height := heightSlice[0]
	height = height[len(`height="`) : len(height)-1]

	// frameRate
	frameRateRe := regexp.MustCompile(`frameRate="[0-9/]*"`)
	frameRateSlice := frameRateRe.FindAllString(videoMPD, -1)
	if len(frameRateSlice) != 1 {
		return nil, fmt.Errorf("frameRate not found")
	}
	frameRate := frameRateSlice[0]
	frameRate = frameRate[len(`frameRate="`) : len(frameRate)-1]

	// bandwidth
	bandwidthRe := regexp.MustCompile(`bandwidth="[0-9]*"`)
	bandwidthSlice := bandwidthRe.FindAllString(videoMPD, -1)
	if len(bandwidthSlice) != 1 {
		return nil, fmt.Errorf("bandwidth not found")
	}
	bandwidth := bandwidthSlice[0]
	bandwidth = bandwidth[len(`bandwidth="`) : len(bandwidth)-1]

	// timescale
	timescaleRe := regexp.MustCompile(`timescale="[0-9]*"`)
	timescaleSlice := timescaleRe.FindAllString(videoMPD, -1)
	if len(timescaleSlice) != 1 {
		return nil, fmt.Errorf("timescaleSlice not found")
	}
	timescale := timescaleSlice[0]
	timescale = timescale[len(`timescale="`) : len(timescale)-1]

	// duration
	durationRe := regexp.MustCompile(`duration="[0-9]*"`)
	durationSlice := durationRe.FindAllString(videoMPD, -1)
	if len(durationSlice) != 1 {
		return nil, fmt.Errorf("duration not found")
	}
	duration := durationSlice[0]
	duration = duration[len(`duration="`) : len(duration)-1]

	va := &videoAttr{
		minBufferTime:             minBufferTime,
		mediaPresentationDuration: mediaPresentationDuration,
		codecs:    codecs,
		par:       par,
		width:     width,
		height:    height,
		frameRate: frameRate,
		bandwidth: bandwidth,
		timescale: timescale,
		duration:  duration,
	}
	return va, nil
}

func getAudioAttr(audioMPD string) (*audioAttr, error) {
	// minBufferTime
	minBufferTimeRe := regexp.MustCompile(`minBufferTime="[PYMDTHMS0-9.]*"`)
	minBufferTimeSlice := minBufferTimeRe.FindAllString(audioMPD, -1)
	if len(minBufferTimeSlice) != 1 {
		return nil, fmt.Errorf("minBufferTimeSlice not found")
	}
	minBufferTime := minBufferTimeSlice[0]
	minBufferTime = minBufferTime[len(`minBufferTime="`) : len(minBufferTime)-1]

	// mediaPresentationDuration
	mediaPresentationDurationRe := regexp.MustCompile(`mediaPresentationDuration="[PYMDTHMS0-9.]*"`)
	mediaPresentationDurationSlice := mediaPresentationDurationRe.FindAllString(audioMPD, -1)
	if len(mediaPresentationDurationSlice) != 1 {
		return nil, fmt.Errorf("mediaPresentationDuration not found")
	}
	mediaPresentationDuration := mediaPresentationDurationSlice[0]
	mediaPresentationDuration = mediaPresentationDuration[len(`mediaPresentationDuration="`) : len(mediaPresentationDuration)-1]

	// codecs
	codecsRe := regexp.MustCompile(`codecs="[A-Za-z0-9._-]*"`)
	codecsSlice := codecsRe.FindAllString(audioMPD, -1)
	if len(codecsSlice) != 1 {
		return nil, fmt.Errorf("codecs not found")
	}
	codecs := codecsSlice[0]
	codecs = codecs[len(`codecs="`) : len(codecs)-1]

	// audioSamplingRate
	audioSamplingRateRe := regexp.MustCompile(`audioSamplingRate="[0-9]*"`)
	audioSamplingRateSlice := audioSamplingRateRe.FindAllString(audioMPD, -1)
	if len(audioSamplingRateSlice) != 1 {
		return nil, fmt.Errorf("audioSamplingRate not found")
	}
	audioSamplingRate := audioSamplingRateSlice[0]
	audioSamplingRate = audioSamplingRate[len(`audioSamplingRate="`) : len(audioSamplingRate)-1]

	// bandwidth
	bandwidthRe := regexp.MustCompile(`bandwidth="[0-9]*"`)
	bandwidthSlice := bandwidthRe.FindAllString(audioMPD, -1)
	if len(bandwidthSlice) != 1 {
		return nil, fmt.Errorf("bandwidth not found")
	}
	bandwidth := bandwidthSlice[0]
	bandwidth = bandwidth[len(`bandwidth="`) : len(bandwidth)-1]

	// timescale
	timescaleRe := regexp.MustCompile(`timescale="[0-9]*"`)
	timescaleSlice := timescaleRe.FindAllString(audioMPD, -1)
	if len(timescaleSlice) != 1 {
		return nil, fmt.Errorf("timescaleSlice not found")
	}
	timescale := timescaleSlice[0]
	timescale = timescale[len(`timescale="`) : len(timescale)-1]

	// duration
	durationRe := regexp.MustCompile(`duration="[0-9]*"`)
	durationSlice := durationRe.FindAllString(audioMPD, -1)
	if len(durationSlice) != 1 {
		return nil, fmt.Errorf("duration not found")
	}
	duration := durationSlice[0]
	duration = duration[len(`duration="`) : len(duration)-1]

	aa := &audioAttr{
		minBufferTime:             minBufferTime,
		mediaPresentationDuration: mediaPresentationDuration,
		codecs:            codecs,
		audioSamplingRate: audioSamplingRate,
		bandwidth:         bandwidth,
		timescale:         timescale,
		duration:          duration,
	}

	return aa, nil
}

func runCmd(name, dir string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if dir != "" {
		cmd.Dir = dir
	}

	outBuf := new(bytes.Buffer)
	cmd.Stdout = outBuf
	errBuf := new(bytes.Buffer)
	cmd.Stderr = errBuf
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ConvertToMP4 converts a non-mp4 video to a mp4 video.
func ConvertToMP4(videoSrcPath, videoDestDir string) error {
	outMP4Path := path.Join(videoDestDir, "out.mp4")
	args := []string{
		// "-noautorotate", // not supported for version smaller 2.7
		"-i", videoSrcPath,
		"-codec:v", "libx264",
		"-profile", "high",
		"-threads", "0",
		"-codec:a", "libfdk_aac",
		"-b:a", "128k",
		outMP4Path,
	}

	err := runCmd(FFMPEG, "", args...)
	if err != nil {
		return nil
	}

	return nil
}

// MPDGen runs MP4Box command to generate a Dash manifest file for the specified
// video at videoSrcPath. It also generates the audio and video segments needed
// by the manifest file.
//
// MP4Box -dash 1000  -dash-profile live -segment-name video -out video.mpd  output.mp4#video
func MPDGen(videoSrcPath, videoDestDir, segment string) error {
	// dash video
	videoArgs := []string{"-dash", segment, "-dash-profile", "live", "-segment-name", "video", "-out"}
	videoMPDPath := path.Join(videoDestDir, "video.mpd")
	videoArgs = append(videoArgs, videoMPDPath)
	videoArgs = append(videoArgs, videoSrcPath+"#video")

	err := runCmd(MP4BOX, videoDestDir, videoArgs...)
	if err != nil {
		return err
	}

	// dash audio
	audioArgs := []string{"-dash", segment, "-dash-profile", "live", "-segment-name", "audio", "-out"}
	audioMPDPath := path.Join(videoDestDir, "audio.mpd")
	audioArgs = append(audioArgs, audioMPDPath)
	audioArgs = append(audioArgs, videoSrcPath+"#audio")

	err = runCmd(MP4BOX, videoDestDir, audioArgs...)
	if err != nil {
		return err
	}

	videoMPD, err := ioutil.ReadFile(videoMPDPath)
	if err != nil {
		return err
	}

	va, err := getVideoAttr(string(videoMPD))
	if err != nil {
		return err
	}

	audioMPD, err := ioutil.ReadFile(audioMPDPath)
	if err != nil {
		return err
	}
	aa, err := getAudioAttr(string(audioMPD))
	if err != nil {
		return err
	}

	manifest := fmt.Sprintf(mpdTpl,
		va.minBufferTime,
		va.mediaPresentationDuration,
		va.codecs,
		va.par,
		va.width,
		va.height,
		va.frameRate,
		va.bandwidth,
		va.timescale,
		va.duration,
		aa.codecs,
		aa.audioSamplingRate,
		aa.bandwidth,
		aa.timescale,
		aa.duration)

	mpdPath := path.Join(videoDestDir, "manifest.mpd")
	return ioutil.WriteFile(mpdPath, []byte(manifest), 0755)
}

// MP3ToMP4 converts a mp3 file to a mp4 file.
func MP3ToMP4(mp3File, mp4File string) error {
	args := []string{
		"-i", mp3File,
		"-vn", mp4File,
	}
	return runCmd(FFMPEG, "", args...)
}

// MP3ToMP4Batch converts mp3 files to mp4 format in batch.
func MP3ToMP4Batch(srcDir, destDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.ToLower(path.Ext(f.Name())) != ".mp3" {
			continue
		}
		mp3Path := path.Join(srcDir, f.Name())
		mp4Path := path.Join(destDir, f.Name()[:len(f.Name())-4]+".mp4")
		// replace ws by "_" in the file name
		mp4Path = strings.Replace(mp4Path, " ", "_", -1)
		err = MP3ToMP4(mp3Path, mp4Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateMP4List(srcDir, listFilePath string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	list := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.ToLower(path.Ext(f.Name())) != ".mp4" {
			continue
		}
		mp4FilePath := path.Join(srcDir, f.Name())
		list = append(list, fmt.Sprintf("file '%s'", mp4FilePath))
	}
	listStr := strings.Join(list, "\n")
	err = ioutil.WriteFile(listFilePath, []byte(listStr), 0777)
	if err != nil {
		return err
	}

	return nil
}

// MPDGenMP3 generates Dash manifest file and audio segments for the specified video file
// (in mp4 foramt) at videoSrcPath.
func MPDGenMP3(videoSrcPath, videoDestDir, segment string) error {
	// dash audio
	audioArgs := []string{"-dash", segment, "-dash-profile", "live", "-segment-name", "audio", "-out"}
	audioMPDPath := path.Join(videoDestDir, "manifest.mpd")
	audioArgs = append(audioArgs, audioMPDPath)
	audioArgs = append(audioArgs, videoSrcPath+"#audio")

	return runCmd(MP4BOX, videoDestDir, audioArgs...)
}

// ConcatenateVideos concatenates multiple mp4 files into a single file.
func ConcatenateVideos(srcDir, destDir string) error {
	listFilePath := path.Join(destDir, "list.txt")
	err := generateMP4List(srcDir, listFilePath)
	if err != nil {
		return err
	}

	outMP4Path := path.Join(destDir, "out.mp4")
	// ffmpeg -f concat -i mylist.txt -c copy output.mp4
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", listFilePath,
		"-c", "copy",
		outMP4Path,
	}
	return runCmd(FFMPEG, "", args...)
}
