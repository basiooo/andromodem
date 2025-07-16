package common_service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/basiooo/andromodem/pkg/adb_processor/processor"
	"github.com/basiooo/andromodem/pkg/adb_processor/utils"
	"github.com/basiooo/andromodem/pkg/cache"
	adb "github.com/basiooo/goadb"
)

func GetAndroidVersion(device *adb.Device, adbProcessor processor.IProcessor, useCache bool) (uint8, error) {
	serial, _ := device.Serial()
	cacheKey := fmt.Sprintf("android_version_%s", serial)
	cacheInstance := cache.GetInstance()
	if useCache {
		if cachedVersion, found := cacheInstance.Get(cacheKey); found {
			return cachedVersion.(uint8), nil
		}
	}
	if androidVersion, err := adbProcessor.Run(device, command.GetAndroidVersionCommand, false); err == nil {
		result := androidVersion.(*parser.RawParser)
		v := strings.Split(result.Result, ".")[0]
		majorVersion, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		version := uint8(majorVersion)
		cacheInstance.Set(cacheKey, version, 5*time.Minute)

		return version, nil
	}
	return 0, nil
}

func GetDeviceRootAndAccessInfo(device *adb.Device, adbProcessor processor.IProcessor, useCache bool) (*model.DeviceRootInfo, error) {
	serial, _ := device.Serial()
	cacheKey := fmt.Sprintf("device_root_info_%s", serial)
	cacheInstance := cache.GetInstance()
	if useCache {
		if cachedInfo, found := cacheInstance.Get(cacheKey); found {
			return cachedInfo.(*model.DeviceRootInfo), nil
		}
	}
	result := &model.DeviceRootInfo{
		RootMethod:  "",
		Rooted:      false,
		ShellAccess: false,
	}

	rootInfo, err := adbProcessor.RunWithRoot(device, command.GetRootCommand)
	if err != nil {
		cacheInstance.Set(cacheKey, result, 5*time.Minute)
		return result, nil
	}
	if root, ok := rootInfo.(*parser.Root); ok {
		result.Rooted = root.IsRooted

		if root.IsRooted {
			shellRootAccess, err := adbProcessor.RunWithRoot(device, command.GetDeviceRootAccessCommand)
			if err == nil {
				result.ShellAccess = strings.Contains(utils.GetResultFromRaw(shellRootAccess), "1")
			}
		}
	}
	cacheInstance.Set(cacheKey, result, 5*time.Minute)
	return result, nil
}
