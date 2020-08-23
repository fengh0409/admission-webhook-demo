package config

import (
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

// Conf is the global config
var Conf *Config

type Config struct {
	Sidecar *SidecarConfig
	Envs    *EnvConfig
}

// SidecarConfig the config object
type SidecarConfig struct {
	Containers []corev1.Container `yaml:"containers"`
	Volumes    []corev1.Volume    `yaml:"volumes"`
}

type EnvConfig struct {
	Mode string
}

// NewConfig init config and add file watcher
func NewConfig(path string) error {
	if err := NewWatcher(path); err != nil {
		return err
	}

	sidecarConfig, err := LoadConfig(path)
	if err != nil {
		return err
	}

	Conf = &Config{
		Sidecar: sidecarConfig,
		Envs:    &EnvConfig{},
	}

	return nil
}

// LoadConfig load the config path
func LoadConfig(path string) (*SidecarConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sidecarConfig *SidecarConfig
	err = yaml.Unmarshal(data, &sidecarConfig)
	if err != nil {
		return nil, err
	}

	glog.Infof("%+v", sidecarConfig)

	return sidecarConfig, nil
}

// NewWatcher watch the file change, and reload it
func NewWatcher(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("newWatcher failed,error:%s", err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
					sidecarConfig, err := LoadConfig(event.Name)
					if err != nil {
						glog.Infof("hot load config error:%s", err)
					}
					Conf.Sidecar = sidecarConfig
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					err := watcher.Add(event.Name)
					if err != nil {
						glog.Infof("remove file event:%s err:%s", event, err)
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					glog.Infof("watcher error:%s", err)
				}
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		return fmt.Errorf("add watcher failed,error:%s", err)
	}

	return nil
}
