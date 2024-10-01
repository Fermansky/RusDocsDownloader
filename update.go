package main

import (
	"context"
	"github.com/google/go-github/v65/github"
	"github.com/inconshreveable/go-update"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/exec"
)

func CheckUpdate() (string, error) {
	client := github.NewClient(http.DefaultClient)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "Fermansky", "RusDocsDownloader")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if GetVersion() != release.GetTagName() {
		return release.GetTagName(), nil
	}

	return "", nil
}

const version = "v0.1.0-SNAPSHOT"

func GetVersion() string {
	return version
}

func restart() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath)
	return cmd.Start()
}

func ApplyUpdate() error {
	client := github.NewClient(http.DefaultClient)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "Fermansky", "RusDocsDownloader")
	if err != nil {
		return errors.WithStack(err)
	}

	if GetVersion() == release.GetTagName() {
		return nil
	}

	var link string
	for _, asset := range release.Assets {
		if asset.GetName() == "RusDocsDownloader.exe" {
			link = asset.GetBrowserDownloadURL()
			break
		}
	}

	if link != "" {
		proxyLink := `https://mirror.ghproxy.com/` + link
		//logger.Logger.Infof("代理链接:%s", proxyLink)
		//优先尝试用国内代理下载
		respProxy, err := http.Get(proxyLink)
		if os.IsTimeout(err) {
			//logger.Logger.Infof("源链接:%s", link)
			resp, err := http.Get(link)
			if err != nil {
				return errors.WithStack(err)
			}
			defer resp.Body.Close()
			//logger.Logger.Infof("下载成功")
			err = update.Apply(respProxy.Body, update.Options{})
			if err != nil {
				return errors.WithStack(err)
			}
			//logger.Logger.Infof("更新成功")

			err = restart()
			if err != nil {
				return errors.WithStack(err)
			}
			//logger.Logger.Infof("重启成功")
			os.Exit(1)
		} else if err != nil {
			return errors.WithStack(err)
		}
		defer respProxy.Body.Close()
		//logger.Logger.Infof("使用代理下载成功")

		err = update.Apply(respProxy.Body, update.Options{})
		if err != nil {
			return errors.WithStack(err)
		}
		//logger.Logger.Infof("更新成功")

		err = restart()
		if err != nil {
			return errors.WithStack(err)
		}
		//logger.Logger.Infof("重启成功")
		os.Exit(1)
	}

	return nil
}
