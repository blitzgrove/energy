//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package internal

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"encoding/json"
	"fmt"
	progressbar "github.com/energye/energy/v2/cmd/internal/progress-bar"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var CmdInstall = &Command{
	UsageLine: "install -p [path] -v [version] -n [name] -d [download]",
	Short:     "Automatically configure the CEF and Energy framework",
	Long: `
	-p Installation directory Default current directory
	-v Specifying a version number,Default latest
	-n Name of the frame after installation
	-d Download Source, gitee or github, Default gitee
	.  Execute default command

Automatically configure the CEF and Energy framework.

During this process, CEF and Energy are downloaded.

Default framework name is "EnergyFramework".
`,
}

type downloadInfo struct {
	fileName      string
	frameworkPath string
	downloadPath  string
	url           string
	success       bool
	isSupport     bool
}

func init() {
	CmdInstall.Run = runInstall
}

const (
	GTK3 = iota + 1
	GTK2
)

// https://cef-builds.spotifycdn.com/cef_binary_107.1.11%2Bg26c0b5e%2Bchromium-107.0.5304.110_windows64.tar.bz2
// 运行安装
func runInstall(c *CommandConfig) error {
	var (
		gtk int
	)
	if isLinux { // linux, gtk2, gtk3 install
		var gtksel = []string{"", "GTK3", "GTK2"}
		println(`
Your current installation environment is Linux and there are two GTK solutions available, namely GTK2 GTK3
	GTK2: Supports LCL window and can use various LCL controls. The CEF<=106 version supports by default
	GTK3: does not support LCL windows and most LCL controls, CEF>=107 version default support
  Enter the number to select the installation support for different versions of GTK framework
  Number Description
	1: GTK3 
	2: GTK2`)
		// 输入选择GTK支持框架
		for {
			print("Please enter the number: ")
			fmt.Scan(&gtk)
			if gtk == GTK3 || gtk == GTK2 {
				break
			} else {
				println("Number input error, please re-enter. Only 1: GTK3 or 2: GTK2 are supported")
			}
		}
		println("Support framework selected:", gtksel[gtk])
	}

	if c.Install.Path == "" {
		c.Install.Path = c.Wd
	}
	installPathName := filepath.Join(c.Install.Path, c.Install.Name)
	println("Install Path", installPathName)
	if c.Install.Version == "" {
		c.Install.Version = "latest"
	}
	os.MkdirAll(c.Install.Path, fs.ModePerm)
	os.MkdirAll(installPathName, fs.ModePerm)
	os.MkdirAll(filepath.Join(c.Install.Path, frameworkCache), fs.ModePerm)
	println("Start downloading CEF and Energy dependency")
	downloadJSON, err := downloadConfig(DownloadVersionConfigUrl)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	var edv map[string]interface{}
	downloadJSON = bytes.TrimPrefix(downloadJSON, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(downloadJSON, &edv); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	var latest = edv["latest"].(string)
	var versionList = edv["versionList"].(map[string]interface{})

	var version map[string]interface{}
	if c.Install.Version == "latest" {
		if v, ok := versionList[latest]; ok {
			version = v.(map[string]interface{})
		}
	} else {
		if c.Install.Version[0] != 'v' {
			c.Install.Version = string('v') + c.Install.Version
		}
		if v, ok := versionList[c.Install.Version]; ok {
			version = v.(map[string]interface{})
		}
	}
	println("Check version")
	if version == nil || len(version) == 0 {
		println("Invalid version number:", c.Install.Version)
		os.Exit(1)
	}
	var (
		versionCEF, versionENERGY string
	)
	if isLinux {
		if gtk == GTK3 {
			versionCEF = ToRNilString(version["cef"], "")
		} else {
			versionCEF = ToRNilString(version["cefgtk2"], "")
		}
	} else {
		versionCEF = ToString(version["cef"])
	}
	versionENERGY = version["energy"].(string)
	var downloadURL map[string]interface{}
	if c.Install.Download == "gitee" {
		downloadURL = edv["gitee"].(map[string]interface{})
	} else if c.Install.Download == "github" {
		downloadURL = edv["github"].(map[string]interface{})
	} else {
		println("Invalid download source, only support github or gitee:", c.Install.Download)
		os.Exit(1)
	}
	libCEFOS, isSupport := cefOS()
	libEnergyOS, isSupport := energyOS(gtk)
	var downloadCefURL = downloadURL["cefURL"].(string)
	var downloadEnergyURL = downloadURL["energyURL"].(string)
	downloadCefURL = strings.ReplaceAll(downloadCefURL, "{version}", versionCEF)
	downloadCefURL = strings.ReplaceAll(downloadCefURL, "{OSARCH}", libCEFOS)
	downloadEnergyURL = strings.ReplaceAll(downloadEnergyURL, "{version}", versionENERGY)
	downloadEnergyURL = strings.ReplaceAll(downloadEnergyURL, "{OSARCH}", libEnergyOS)

	//提取文件配置
	extractData, err := downloadConfig(DownloadExtractUrl)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error(), "\n")
		os.Exit(1)
	}
	var extractConfig map[string]interface{}
	extractData = bytes.TrimPrefix(extractData, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(extractData, &extractConfig); err != nil {
		fmt.Fprint(os.Stderr, err.Error(), "\n")
		os.Exit(1)
	}
	extractOSConfig := extractConfig[runtime.GOOS].(map[string]interface{})

	var downloads = make(map[string]*downloadInfo)
	downloads[cefKey] = &downloadInfo{isSupport: isSupport, fileName: urlName(downloadCefURL), downloadPath: filepath.Join(c.Install.Path, frameworkCache, urlName(downloadCefURL)), frameworkPath: installPathName, url: downloadCefURL}
	downloads[energyKey] = &downloadInfo{isSupport: isSupport, fileName: urlName(downloadEnergyURL), downloadPath: filepath.Join(c.Install.Path, frameworkCache, urlName(downloadEnergyURL)), frameworkPath: installPathName, url: downloadEnergyURL}
	for key, dl := range downloads {
		fmt.Printf("Download %s: %s\n", key, dl.url)
		if !dl.isSupport {
			println("energy command line does not support the system architecture download 【", dl.fileName, "]")
			continue
		}
		bar := progressbar.NewBar(100)
		bar.SetNotice("\t")
		bar.HideRatio()
		err = downloadFile(dl.url, dl.downloadPath, func(totalLength, processLength int64) {
			bar.PrintBar(int((float64(processLength) / float64(totalLength)) * 100))
		})
		bar.PrintEnd("Download " + dl.fileName + " success")
		if err != nil {
			println("Download", dl.fileName, "error", err)
		}
		dl.success = err == nil
	}
	println("Unpack files")
	var removeFileList = make([]string, 0, 0)
	for key, di := range downloads {
		if !di.isSupport {
			println("energy command line does not support the system architecture, continue.")
			continue
		}
		if di.success {
			if key == cefKey {
				bar := progressbar.NewBar(0)
				bar.SetNotice("Unpack file " + key + ": ")
				tarName := UnBz2ToTar(di.downloadPath, func(totalLength, processLength int64) {
					bar.PrintSizeBar(processLength)
				})
				bar.PrintEnd()
				ExtractFiles(key, tarName, di, extractOSConfig)
				removeFileList = append(removeFileList, tarName)
			} else if key == energyKey {
				ExtractFiles(key, di.downloadPath, di, extractOSConfig)
			}
			println("Unpack file", key, "success\n")
		}
	}
	for _, rmFile := range removeFileList {
		println("Remove file", rmFile)
		os.Remove(rmFile)
	}
	setEnergyHomeEnv(EnergyHomeKey, installPathName)
	println("\n", CmdInstall.Short, "SUCCESS \nVersion:", c.Install.Version, "=>", versionENERGY)
	return nil
}

func cefOS() (string, bool) {
	if isWindows { // windows arm for 64 bit, windows for 32/64 bit
		if runtime.GOARCH == "arm64" {
			return "windowsarm64", true
		}
		return fmt.Sprintf("windows%d", strconv.IntSize), true
	} else if isLinux { //linux for 64 bit
		if runtime.GOARCH == "arm64" {
			return "linuxarm64", true
		} else if runtime.GOARCH == "amd64" {
			return "linux64", true
		}
	} else if isDarwin { // macosx for 64 bit
		if runtime.GOARCH == "arm64" {
			return "macosarm64", true
		} else if runtime.GOARCH == "amd64" {
			return "macosx64", true
		}
	}
	//not support
	return fmt.Sprintf("%v %v", runtime.GOOS, runtime.GOARCH), false
}

func energyOS(gtk int) (string, bool) {
	if isWindows {
		return fmt.Sprintf("Windows %d bits", strconv.IntSize), true
	} else if isLinux {
		if gtk == GTK3 {
			return "Linux x86 64 bits", true
		}
		return "Linux GTK2 x86 64 bits", true
	} else if isDarwin {
		return "MacOSX x86 64 bits", true
	}
	//not support
	return fmt.Sprintf("%v %v", runtime.GOOS, runtime.GOARCH), false
}

// 提取文件
func ExtractFiles(keyName, sourcePath string, di *downloadInfo, extractOSConfig map[string]interface{}) {
	println("Extract", keyName, "sourcePath:", sourcePath, "targetPath:", di.frameworkPath)
	files := extractOSConfig[keyName].([]interface{})
	if keyName == cefKey {
		//tar
		ExtractUnTar(sourcePath, di.frameworkPath, files...)
	} else if keyName == energyKey {
		//zip
		ExtractUnZip(sourcePath, di.frameworkPath, files...)
	}
}

func filePathInclude(compressPath string, files ...interface{}) (string, bool) {
	for _, file := range files {
		f := file.(string)
		tIdx := strings.LastIndex(f, "/*")
		if tIdx != -1 {
			f = f[:tIdx]
			if f[0] == '/' {
				if strings.Index(compressPath, f[1:]) == 0 {
					return compressPath, true
				}
			}
			if strings.Index(compressPath, f) == 0 {
				return strings.Replace(compressPath, f, "", 1), true
			}
		} else {
			if f[0] == '/' {
				if compressPath == f[1:] {
					return f, true
				}
			}
			if compressPath == f {
				f = f[strings.Index(f, "/")+1:]
				return f, true
			}
		}
	}
	return "", false
}

func dir(path string) string {
	path = strings.ReplaceAll(path, "\\", string(filepath.Separator))
	lastSep := strings.LastIndex(path, string(filepath.Separator))
	return path[:lastSep]
}

func ExtractUnTar(filePath, targetPath string, files ...interface{}) {
	reader, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: cannot read tar file, error=[%v]\n", err)
		return
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error: cannot read tar file, error=[%v]\n", err)
			os.Exit(1)
			return
		}
		compressPath := header.Name[strings.Index(header.Name, "/")+1:]
		includePath, isInclude := filePathInclude(compressPath, files...)
		if !isInclude {
			continue
		}
		info := header.FileInfo()
		targetFile := filepath.Join(targetPath, includePath)
		fmt.Println("compressPath:", compressPath, "-> targetFile:", targetFile)
		if info.IsDir() {
			if err = os.MkdirAll(targetFile, info.Mode()); err != nil {
				fmt.Printf("error: cannot mkdir file, error=[%v]\n", err)
				os.Exit(1)
				return
			}
		} else {
			fDir := dir(targetFile)
			_, err = os.Stat(fDir)
			if os.IsNotExist(err) {
				if err = os.MkdirAll(fDir, info.Mode()); err != nil {
					fmt.Printf("error: cannot file mkdir file, error=[%v]\n", err)
					os.Exit(1)
					return
				}
			}
			file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				fmt.Printf("error: cannot open file, error=[%v]\n", err)
				os.Exit(1)
				return
			}
			defer file.Close()
			bar := progressbar.NewBar(100)
			bar.SetNotice("\t")
			bar.HideRatio()
			writeFile(tarReader, file, header.Size, func(totalLength, processLength int64) {
				bar.PrintBar(int((float64(processLength) / float64(totalLength)) * 100))
			})
			file.Sync()
			bar.PrintBar(100)
			bar.PrintEnd()
			if err != nil {
				fmt.Printf("error: cannot write file, error=[%v]\n", err)
				os.Exit(1)
				return
			}
		}
	}
}

func ExtractUnZip(filePath, targetPath string, files ...interface{}) {
	if rc, err := zip.OpenReader(filePath); err == nil {
		defer rc.Close()
		for i := 0; i < len(files); i++ {
			if f, err := rc.Open(files[i].(string)); err == nil {
				defer f.Close()
				st, _ := f.Stat()
				targetFileName := filepath.Join(targetPath, st.Name())
				if st.IsDir() {
					os.MkdirAll(targetFileName, st.Mode())
					continue
				}
				if targetFile, err := os.Create(targetFileName); err == nil {
					fmt.Println("extract file: ", st.Name())
					bar := progressbar.NewBar(100)
					bar.SetNotice("\t")
					bar.HideRatio()
					writeFile(f, targetFile, st.Size(), func(totalLength, processLength int64) {
						bar.PrintBar(int((float64(processLength) / float64(totalLength)) * 100))
					})
					bar.PrintBar(100)
					bar.PrintEnd()
					targetFile.Close()
				}
			} else {
				fmt.Printf("error: cannot open file, error=[%v]\n", err)
				os.Exit(1)
				return
			}
		}
	} else {
		if err != nil {
			fmt.Printf("error: cannot read zip file, error=[%v]\n", err)
			os.Exit(1)
		}
	}
}

// 释放bz2文件到tar
func UnBz2ToTar(name string, callback func(totalLength, processLength int64)) string {
	fileBz2, err := os.Open(name)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		os.Exit(1)
	}
	defer fileBz2.Close()
	dirName := fileBz2.Name()
	dirName = dirName[:strings.LastIndex(dirName, ".")]
	r := bzip2.NewReader(fileBz2)
	w, err := os.Create(dirName)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		os.Exit(1)
	}
	defer w.Close()
	writeFile(r, w, 0, callback)
	return dirName
}

func writeFile(r io.Reader, w *os.File, totalLength int64, callback func(totalLength, processLength int64)) {
	buf := make([]byte, 1024*10)
	var written int64
	for {
		nr, err := r.Read(buf)
		if nr > 0 {
			nw, err := w.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			callback(totalLength, written)
			if err != nil {
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if err != nil {
			break
		}
	}
}

// url文件名
func urlName(downloadUrl string) string {
	if u, err := url.QueryUnescape(downloadUrl); err != nil {
		return ""
	} else {
		u = u[strings.LastIndex(u, "/")+1:]
		return u
	}
}

// 下载文件配置
func downloadConfig(url string) ([]byte, error) {
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	os.Remove(filename)
	return false
}

// 下载文件
func downloadFile(url string, localPath string, callback func(totalLength, processLength int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 1024*10)
		written int64
	)
	tmpFilePath := localPath + ".download"
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		os.Exit(1)
		return err
	}
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		os.Exit(1)
		return err
	}
	if isFileExist(localPath, fsize) {
		return nil
	}
	println("Save path: [", localPath, "] file size:", fsize)
	file, err := os.Create(tmpFilePath)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		os.Exit(1)
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		fmt.Printf("Download-error=[body is null]\n")
		os.Exit(1)
		return nil
	}
	defer resp.Body.Close()
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, err := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			callback(fsize, written)
			if err != nil {
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		file.Sync()
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		if err != nil {
			return err
		}
	}
	return err
}
