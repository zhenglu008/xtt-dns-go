package xtt

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const WINDOWS_HOST = "C:\\Windows\\System32\\drivers\\etc\\hosts"
const LINUX_HOST = "/etc/hosts"

type xttConfig struct {
	HostRemoteAddress  string `yaml:"host_remote_address"`
	TempFilePath       string `yaml:"temp_file_path"`
	LocalHostFile      string `yaml:"local_host_file"`
	RestartExePath     string `yaml:"restart_exe_path"`
	RestartProcessName string `yaml:"restart_process_name"`
	RefreshSecond      time.Duration  `yaml:"refresh_second"`
	LogFile            string `yaml:"log_file"`
}

type XttDns struct {
	Config xttConfig `yaml:"config"`
	RemoteHost string
	LocalHost string
}

// 执行入口
func (xtt XttDns) Run() {

	f, err := os.OpenFile(xtt.Config.LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	timeTickerChan := time.Tick(time.Second * xtt.Config.RefreshSecond)
	for {

		xtt.Print("pull remote host")

		// 本地和远程不同时更新host
		if xtt.remoteHost() != xtt.localHosts() {
			xtt.refresh()
			xtt.restartMiner()
		}

		<-timeTickerChan
	}
}

func (xtt XttDns) restartMiner() {

	if runtime.GOOS == "linux" {
		// TODO
	} else if runtime.GOOS == "windows" {
		xtt.runCmd("ipconfig", []string{"/flushdns"})
		xtt.runCmd("taskkill", []string{"/f", "/im", xtt.Config.RestartProcessName})
		xtt.Print("restart " + xtt.Config.RestartProcessName)
		err := exec.Command(xtt.Config.RestartExePath).Start(); if err != nil {
			xtt.Print("exec error=" + err.Error())
		}
	}
}

func (xtt XttDns) refresh() {

	remoteHosts := xtt.remoteHost()
	xtt.Print("refresh")

	for _, v := range strings.Split(remoteHosts, "\n") {
		xtt.Print(v)
	}

	if runtime.GOOS == "linux" {
		xtt.runCmd("truncate", []string{"-s", "0", LINUX_HOST}) // 清空hosts文件 直接删除文件会报错
		err := ioutil.WriteFile(LINUX_HOST, []byte(remoteHosts), 0666)
		if err != nil {
			xtt.Print("Host写入失败=%v" + err.Error())
		}
	} else if runtime.GOOS == "windows" {
		err := ioutil.WriteFile(xtt.Config.TempFilePath, []byte(remoteHosts), 0666)
		if err != nil {
			xtt.Print("Host写入失败=%v" + err.Error())
		}
		err = MoveFile(xtt.Config.TempFilePath, WINDOWS_HOST); if err != nil {
			xtt.Print("MoveFile error=" + err.Error())
		}
	}

}


// 获取host配置
func (xtt XttDns) remoteHost () string {

	resp, err := http.Get(xtt.Config.HostRemoteAddress)
	if err != nil {
		xtt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body); if err != nil {
		xtt.Println("ioutil.ReadAll error = ", err)
	}

	return string(body)
}

func (xtt XttDns) localHosts() string {

	var data []byte
	var err error
	if runtime.GOOS == "windows" {
		data, err = ioutil.ReadFile(WINDOWS_HOST); if err != nil {
			xtt.Println("ioutil.ReadFile error =", err)
		}
	} else if runtime.GOOS == "linux" {
		data, err = ioutil.ReadFile(LINUX_HOST); if err != nil {
			xtt.Println("ioutil.ReadFile error =", err)
		}
	} else {
		panic(runtime.GOOS + " system does not support")
	}

	return string(data)
}

func (xtt XttDns) runCmd(name string, args []string) string {
	var res []byte
	var err error

	xtt.Print(name + " " + strings.Join(args, " "))

	// 执行单个shell命令时, 直接运行即可
	err = exec.Command(name, args...).Run()
	if err != nil {
		xtt.Print(err.Error())
	}
	// 默认输出有一个换行
	return string(res)
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (xtt XttDns) Print(message string) {
	fmt.Printf("[%s] %s \n", time.Now().Format("2006-01-02 15:04:05"), message)
	log.Printf("[%s] %s \n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func (xtt XttDns) Printf(format string, a...interface{}) {
	date := fmt.Sprintf("[%s] %s \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(date + format, a)
	log.Printf(date + format, a)
}

func (xtt XttDns) Println(a...interface{}) {
	fmt.Println(fmt.Sprintf("[%s] %s \n", time.Now().Format("2006-01-02 15:04:05")), a)
	log.Println(fmt.Sprintf("[%s] %s \n", time.Now().Format("2006-01-02 15:04:05")), a)
}