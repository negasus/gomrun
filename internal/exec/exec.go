package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"

	"github.com/negasus/gomrun/internal/config"
	"github.com/negasus/gomrun/internal/stdlog"
)

func Exec(wg *sync.WaitGroup, cfg *config.Config, name string, nameMaxLen int, serviceColor color.Attribute) {
	defer wg.Done()

	srv, ok := cfg.Services[name]
	if !ok {
		color.Red("unexpected error, service %q not found\n", name)
		return
	}

	if srv.Delay == 0 {
		color.Cyan("starting %q\n", name)
	} else {
		color.Cyan("starting %q with delay %d sec\n", name, srv.Delay)
		time.Sleep(time.Duration(srv.Delay) * time.Second)
	}

	target := srv.Cmd

	if srv.Build != nil {
		target = path.Join(os.TempDir(), "gomrun", name+strconv.Itoa(int(time.Now().UnixNano())))

		errBuild := build(target, srv)
		if errBuild != nil {
			color.Red("error build service %q, %v\n", name, errBuild)
			return
		}
	}

	cmd := exec.Command(target, srv.Args...)
	cmd.Dir = srv.WorkDir
	cmd.Stdout = stdlog.New(name, nameMaxLen, serviceColor)
	cmd.Stderr = stdlog.New(name, nameMaxLen, serviceColor)
	for _, envsetName := range srv.Envset {
		es, ok := cfg.Envset[envsetName]
		if !ok {
			color.Red("envset %q not found\n", envsetName)
			return
		}
		for n, v := range es {
			cmd.Env = append(cmd.Env, n+"="+v)
		}
	}
	for n, v := range srv.Environment {
		cmd.Env = append(cmd.Env, n+"="+v)
	}

	if srv.EnvFile != "" {
		errEnvFile := loadEnvFile(cmd, srv.EnvFile)
		if errEnvFile != nil {
			color.Red("error load env file %q, %v\n", srv.EnvFile, errEnvFile)
			return
		}
	}

	errRun := cmd.Start()
	if errRun != nil {
		color.Red("service %q start error, %v\n", name, errRun)
		return
	}

	errWait := cmd.Wait()
	if errWait != nil {
		color.Red("error wait service stop %q, %v\n", name, errWait)
		return
	}

	color.Cyan("service %q terminated\n", name)
}

func loadEnvFile(cmd *exec.Cmd, file string) error {
	data, errReadFile := os.ReadFile(file)
	if errReadFile != nil {
		return fmt.Errorf("error read env file %q, %w", file, errReadFile)
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.Contains(line, "=") {
			return fmt.Errorf("invalid env line %q", line)
		}
		cmd.Env = append(cmd.Env, line)
	}

	return nil
}

func build(target string, srv config.Service) error {
	cmd := exec.Command("go", "build", "-o", target, srv.Build.Path)
	cmd.Dir = srv.Build.Context
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
