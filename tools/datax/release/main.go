// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	sourceUserPath = "datax/"
	destUserPath   = "release/datax/"

	sourceExamplePath = "cmd/datax/examples"
	destExamplePath   = "release/examples"
)

func main() {
	err := copyMarkdown("plugin/reader")
	if err != nil {
		fmt.Println("copyMarkdown reader fail. err:", err)
		os.Exit(1)
	}

	err = copyMarkdown("plugin/writer")
	if err != nil {
		fmt.Println("copyMarkdown writer fail. err:", err)
		os.Exit(1)
	}

	err = copyConfig()
	if err != nil {
		fmt.Println("copyConfig fail. err:", err)
		os.Exit(1)
	}

	data, err := os.ReadFile("README_USER.md")
	if err != nil {
		fmt.Println("ReadFile README_USER.md fail. err:", err)
		os.Exit(1)
	}

	err = os.WriteFile("release/README_USER.md", data, os.ModePerm)
	if err != nil {
		fmt.Println("WriteFile release/README_USER.md fail. err:", err)
		os.Exit(1)
	}

	data, err = os.ReadFile("README_USER_zh-CN.md")
	if err != nil {
		fmt.Println("ReadFile README_USER.md fail. err:", err)
		os.Exit(1)
	}

	err = os.WriteFile("release/README_USER_zh-CN.md", data, os.ModePerm)
	if err != nil {
		fmt.Println("WriteFile release/README_USER.md fail. err:", err)
		os.Exit(1)
	}

	output := ""
	if output, err = cmdOutput("git", "describe", "--abbrev=0", "--tags"); err != nil {
		fmt.Printf("use git to tag version fail. error: %v\n", err)
		os.Exit(1)
	}
	tagVersion := strings.ReplaceAll(output, "\r", "")
	tagVersion = strings.ReplaceAll(tagVersion, "\n", "")
	os.MkdirAll("release/bin", os.ModePerm)
	if runtime.GOOS == "windows" {
		os.Rename("cmd/datax/datax.exe", "release/bin/go-etl.exe")
		if err = zipDir("release", "go-etl-"+tagVersion+"-windows-x86_64.zip"); err != nil {
			fmt.Printf("uzipDir fail. error: %v\n", err)
			os.Exit(1)
		}
	} else if runtime.GOOS == "linux" {
		os.Rename("cmd/datax/datax", "release/bin/go-etl")
		if err = tarDir("release", "go-etl-"+tagVersion+"-linux-x86_64.tar.gz"); err != nil {
			fmt.Printf("tarDir fail. error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("OS: %v\n", runtime.GOOS)
		os.Exit(1)
	}
}

func copyMarkdown(path string) (err error) {
	var list []fs.DirEntry
	list, err = os.ReadDir(filepath.Join(sourceUserPath, path))
	if err != nil {
		return err
	}
	var data []byte
	for _, v := range list {
		if v.IsDir() {
			data, err = os.ReadFile(filepath.Join(sourceUserPath, path, v.Name(), "README.md"))
			if err != nil {
				err = nil
				continue
			}
			os.MkdirAll(filepath.Join(destUserPath, path, v.Name()), os.ModePerm)
			err = os.WriteFile(filepath.Join(destUserPath, path, v.Name(), "README.md"), data, 0644)
			if err != nil {
				return
			}

			data, err = os.ReadFile(filepath.Join(sourceUserPath, path, v.Name(), "README_zh-CN.md"))
			if err != nil {
				err = nil
				continue
			}
			os.MkdirAll(filepath.Join(destUserPath, path, v.Name()), os.ModePerm)
			err = os.WriteFile(filepath.Join(destUserPath, path, v.Name(), "README_zh-CN.md"), data, 0644)
			if err != nil {
				return
			}
		}
	}
	return
}

func copyConfig() (err error) {
	os.MkdirAll(destExamplePath, os.ModePerm)
	var list []fs.DirEntry
	list, err = os.ReadDir(sourceExamplePath)
	if err != nil {
		return err
	}
	var data []byte
	for _, v := range list {
		if v.IsDir() {
			data, err = os.ReadFile(filepath.Join(sourceExamplePath, v.Name(), "config.json"))
			if err != nil {
				err = nil
				continue
			}
			os.MkdirAll(filepath.Join(destExamplePath, v.Name()), os.ModePerm)
			err = os.WriteFile(filepath.Join(destExamplePath, v.Name(), "config.json"), data, 0644)
			if err != nil {
				return
			}
		}
	}
	return
}

func zipDir(src, dest string) error {
	zipfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

func tarDir(src, dst string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()
	return filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		hdr.Name = strings.TrimPrefix(fileName, string(filepath.Separator))

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		return nil
	})
}

func cmdOutput(cmd string, arg ...string) (output string, err error) {
	c := exec.Command(cmd, arg...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err = c.Run(); err != nil {
		err = fmt.Errorf("%v(%s)", err, stderr.String())
		return
	}
	return stdout.String(), nil
}
