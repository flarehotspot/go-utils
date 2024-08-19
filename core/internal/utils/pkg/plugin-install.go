package pkg

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"core/internal/utils/encdisk"
	"core/internal/utils/git"
	sdkplugin "sdk/api/plugin"
	sdkfs "sdk/utils/fs"
	sdkpaths "sdk/utils/paths"
	sdkstr "sdk/utils/strings"
)

type PluginMetadata struct {
	Def PluginSrcDef
}

type PluginFile struct {
	File     string
	Optional bool
}

var PLuginFiles = []PluginFile{
	{
		File:     "plugin.json",
		Optional: false,
	},
	{
		File:     "plugin.so",
		Optional: false,
	},
	{
		File:     "resources",
		Optional: true,
	},
	{
		File:     "go.mod",
		Optional: false,
	},
	{
		File:     "LICENSE.txt",
		Optional: true,
	},
}

func InstallSrcDef(w io.Writer, def PluginSrcDef) (sdkplugin.PluginInfo, error) {
	switch def.Src {
	case PluginSrcGit:
		return InstallFromGitSrc(w, def)
	case PluginSrcLocal, PluginSrcSystem:
		return InstallFromLocalPath(w, def)
	default:
		return sdkplugin.PluginInfo{}, errors.New("Invalid plugin source: " + def.Src)
	}
}

func InstallFromLocalPath(w io.Writer, def PluginSrcDef) (sdkplugin.PluginInfo, error) {
	w.Write([]byte("Building plugin from local path: " + def.LocalPath))

	info, err := GetSrcInfo(def.LocalPath)
	if err != nil {
		return sdkplugin.PluginInfo{}, err
	}

	err = InstallPlugin(def.LocalPath, InstallOpts{RemoveSrc: false})
	if err != nil {
		return sdkplugin.PluginInfo{}, err
	}

	if err := MarkPluginAsInstalled(def, GetInstallPath(info.Package)); err != nil {
		return sdkplugin.PluginInfo{}, err
	}

	return info, nil
}

func InstallFromGitSrc(w io.Writer, def PluginSrcDef) (sdkplugin.PluginInfo, error) {
	randomPath := RandomPluginPath()
	diskfile := filepath.Join(randomPath, "disk")
	mountpath := filepath.Join(randomPath, "mount")
	clonePath := filepath.Join(mountpath, "clone", "0") // need extra sub dir

	dev := sdkstr.Rand(8)
	mnt := encdisk.NewEncrypedDisk(diskfile, mountpath, dev)
	if err := mnt.Mount(); err != nil {
		log.Println("Error mounting disk: ", err)
		return sdkplugin.PluginInfo{}, err
	}

	defer mnt.Unmount()

	repo := git.RepoSource{URL: def.GitURL, Ref: def.GitRef}

	log.Println("Cloning plugin from git: " + def.GitURL)
	if err := git.Clone(w, repo, clonePath); err != nil {
		log.Println("Error cloning: ", err)
		return sdkplugin.PluginInfo{}, err
	}

	info, err := GetSrcInfo(clonePath)
	if err != nil {
		log.Println("Error getting plugin info: ", err)
		return sdkplugin.PluginInfo{}, err
	}

	if err := InstallPlugin(clonePath, InstallOpts{RemoveSrc: false}); err != nil {
		return sdkplugin.PluginInfo{}, err
	}

	if err := MarkPluginAsInstalled(def, GetInstallPath(info.Package)); err != nil {
		return sdkplugin.PluginInfo{}, err
	}

	return info, nil
}

func InstallPlugin(src string, opts InstallOpts) error {
	log.Println("Installing plugin: ", src)

	var buildpath string

	if opts.Encrypt {
		dev := sdkstr.Rand(8)
		parentPath := RandomPluginPath()
		diskfile := filepath.Join(parentPath, "disk")
		mountpath := filepath.Join(parentPath, "mount")
		buildpath = filepath.Join(mountpath, "build")
		mnt := encdisk.NewEncrypedDisk(diskfile, mountpath, dev)
		if err := mnt.Mount(); err != nil {
			log.Println("Error mounting: ", err)
			return err
		}

		defer mnt.Unmount()
	} else {
		parentpath := filepath.Join(sdkpaths.TmpDir, "plugins", "build", sdkstr.Rand(16))
		buildpath = filepath.Join(parentpath, "0")
		if err := sdkfs.EmptyDir(buildpath); err != nil {
			return err
		}
	}

	if err := BuildPlugin(src, buildpath); err != nil {
		log.Println("Error building plugin: ", err)
		return err
	}

	info, err := GetSrcInfo(src)
	if err != nil {
		return err
	}

	installPath := GetInstallPath(info.Package)
	for _, f := range PLuginFiles {
		err := sdkfs.Copy(filepath.Join(src, f.File), filepath.Join(installPath, f.File))
		if err != nil && !f.Optional {
			return err
		}
	}

	if opts.RemoveSrc {
		if err := os.RemoveAll(src); err != nil {
			return err
		}
	}

	return nil
}
