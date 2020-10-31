package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type repository struct {
	worktree string
	gitDir   string
	viper    *viper.Viper
}

func newRepository(path string, force bool) (*repository, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s.\n\terror is caused by %w", path, err)
	}
	if !force && !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a git repository", path)
	}

	gitDir := fmt.Sprintf("%s/.git/", filepath.Dir(path))
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(gitDir)
	r := &repository{
		worktree: path,
		gitDir:   gitDir,
		viper:    viper.GetViper(),
	}

	_, err = os.Stat(r.gitDir + "config")
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to file.\n\terror is caused by %w", err)
	}
	if err != nil && os.IsNotExist(err) && !force {
		return nil, fmt.Errorf("config file missing")
	}

	if !force {
		err := r.viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
		version, ok := r.viper.GetStringMap("core")["repositoryformatversion"]
		if !ok {
			return nil, fmt.Errorf("repositoryformatversion is not found")
		}
		_version, ok := version.(int)
		if !ok {
			return nil, fmt.Errorf("repositoryformatversion is not integer")
		}
		if _version != 0 {
			return nil, fmt.Errorf("unsupported repositoryformatversion %v", _version)
		}
	}

	return r, nil
}

func (r *repository) createDirIfNeed(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat %s.\n\terror is caused by %w", path, err)
		}
		err = os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to make directory path %s.\n\terror is caused by %w", path, err)
		}
		err = os.Chmod(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to chmod directory path %s to %o.\n\terror is caused by %w", path, 0777, err)
		}
	} else {
		if !fileInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}
	}

	return nil
}

func (r *repository) createFileWithWrite(fileName, content string) error {
	filePath := r.gitDir + fileName
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file path %s.\n\terror is caused by %w", filePath, err)
	}
	n, err := file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write %s", filePath)
	}
	if n == 0 {
		return fmt.Errorf("failed to write %s. write number of bytes is 0", filePath)
	}
	return nil
}

func (r *repository) create() error {
	if err := r.createDirIfNeed(r.worktree); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", r.worktree, err)
	}
	files, err := ioutil.ReadDir(r.worktree)
	if err != nil {
		return fmt.Errorf("failed to read directory path %s.\n\terror is caused by %w", r.worktree, err)
	}
	if len(files) != 0 {
		return fmt.Errorf("%s is not empty", r.worktree)
	}

	if err := r.createDirIfNeed(r.gitDir); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", r.worktree, err)
	}

	toMakeDirs := []string{"branches", "objects", "refs", "refs/tags", "refs/heads"}
	for _, toMakeDir := range toMakeDirs {
		if err = r.createDirIfNeed(r.gitDir + toMakeDir); err != nil {
			return fmt.Errorf("failed to dirs.\n\terror is caused by %w", err)
		}
	}

	toWriteFileMap := map[string]string{
		"description": "unnamed repository: edit this file 'description' to name the repository.",
		"HEAD":        "ref: refs/heads/master\n",
		"config": `[core]
		repositoryformatversion = 0
		filemode = false
		bare = false
	`,
	}
	for fileName, content := range toWriteFileMap {
		if err = r.createFileWithWrite(fileName, content); err != nil {
			return fmt.Errorf("failed to create and write.\n\terror is caused by %w", err)
		}
	}

	return nil
}
