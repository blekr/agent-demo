package biz

import (
	"fmt"
	"github.com/blekr/agent-demo/errors"
	. "github.com/blekr/agent-demo/util"
	procUtil "github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const maxTry = 3

type Command struct {
	Path string
	Args []string
}

type ProcessStat struct {
	Pid int
	IsRunning bool
	CPUPercent float64
	NumConnections int
	CreateTime int64
	MemoryPercent float32
	NumThreads int
}

type Null struct {

}

func (null *Null) Read(p []byte) (int, error) {
	return 0, nil
}
func (null *Null) Write(p []byte) (int, error) {
	return len(p), nil
}

func StartProcess(command *Command) (string, error) {
	ILog.Printf("StartProcess biz: %v", command)

	for retry := 1;; {
		cmd := exec.Command(command.Path, command.Args...)

		// prevent child process from being killed when parent process exit abnormally
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    0,
		}

		err := cmd.Start()
		if err == nil {
			return fmt.Sprintf("%d", cmd.Process.Pid), nil
		}
		if retry < maxTry {
			ILog.Printf("retry starting command (%v): %v", retry, command.Path)
			time.Sleep(3 * time.Second)
			retry ++
		} else {
			return "", &errors.AppError{Code: "MAX_TRY_ERROR", Message: err.Error()}
		}
	}
}

// Graceful shutdown the process:
// Send terminate signal to the process and then wait for it to terminate.
// If it does not terminate in 10 seconds, kill it immediately.
func StopProcess(pid int) error {
	ILog.Printf("StopProcess biz: %v", pid)

	exist, err := procUtil.PidExists(int32(pid))
	if err != nil {
		return err
	}

	if !exist {
		return &errors.AppError{Code: "NOT_EXIST_ERROR", Message: fmt.Sprintf("pid %v not exists", pid)}
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	timeout := time.After(10 * time.Second)
	exited := make(chan error)

	go func() {
		_, err := process.Wait()
		exited <- err
	}()

	select {
	case <-timeout:
		_ = process.Signal(syscall.SIGKILL)
		ILog.Printf("timedout waiting for %v to terminate", pid)
		return nil
	case err := <-exited:
		ILog.Printf("process.wait returned (%v)", err)
		return err
	}
}

func ShowProcess(pid int) (*ProcessStat, error) {
	ILog.Printf("ShowProcess biz: %v", pid)

	exist, err := procUtil.PidExists(int32(pid))
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, &errors.AppError{Code: "NOT_EXIST_ERROR", Message: fmt.Sprintf("pid %v not exists", pid)}
	}

	process, err := procUtil.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	isRunning, _ := process.IsRunning()
	cpuPercent, _ := process.CPUPercent()
	connections, _ := process.Connections()
	createTime, _ := process.CreateTime()
	memoryPercent, _ := process.MemoryPercent()
	numThreads, _ := process.NumThreads()

	return &ProcessStat{
		Pid: pid,
		IsRunning: isRunning,
		CPUPercent: cpuPercent,
		NumConnections: len(connections),
		CreateTime: createTime,
		MemoryPercent: memoryPercent,
		NumThreads: int(numThreads),
	}, nil

}