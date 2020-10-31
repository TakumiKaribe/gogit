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

	configFilePath, err := r.file("config", false)
	_, err = os.Stat(configFilePath)
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

func (r *repository) path(path string) string {
	return r.gitDir + path
}

func (r *repository) file(path string, mkdir bool) (string, error) {
	_, err := r.dir(path, mkdir)
	if err != nil {
		return "", err
	}
	return r.path(path), nil
}

func (r *repository) dir(path string, mkdir bool) (string, error) {
	path = r.path(path)
	fileInfo, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat %s.\n\terror is caused by %w", path, err)
	}

	if err != nil && os.IsNotExist(err) {
		if mkdir {
			err = os.Mkdir(path, 0755)
			if err != nil {
				return "", fmt.Errorf("failed to make directory path %s.\n\terror is caused by %w", path, err)
			}
			err = os.Chmod(path, 0755)
			if err != nil {
				return "", fmt.Errorf("failed to chmod directory path %s to %o.\n\terror is caused by %w", path, 0777, err)
			}
			return path, nil
		}
		return "", nil
	}

	if fileInfo.IsDir() {
		return path, nil
	}

	return "", fmt.Errorf("%s is not a directory", path)
}

func (r *repository) create() error {
	fileInfo, err := os.Stat(r.worktree)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s.\n\terror is caused by %w", r.worktree, err)
	}
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(r.worktree, 0755)
		if err != nil {
			return fmt.Errorf("failed to make directory path %s.\n\terror is caused by %w", r.worktree, err)
		}
		err = os.Chmod(r.gitDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to chmod directory path %s to %o.\n\terror is caused by %w", r.worktree, 0777, err)
		}
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", r.worktree)
	}
	files, err := ioutil.ReadDir(r.worktree)
	if err != nil {
		return fmt.Errorf("failed to read directory path %s.\n\terror is caused by %w", r.worktree, err)
	}
	if len(files) != 0 {
		return fmt.Errorf("%s is not empty", r.worktree)
	}

	fileInfo, err = os.Stat(r.gitDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat path %s.\n\terror is caused by %w", r.gitDir, err)
	}
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(r.gitDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to make directory path %s.\n\terror is caused by %w", r.gitDir, err)
		}
		err = os.Chmod(r.gitDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to chmod directory path %s to %o.\n\terror is caused by %w", r.gitDir, 0777, err)
		}
	}

	if _, err = r.dir("branches", true); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", "branches", err)
	}
	if _, err = r.dir("objects", true); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", "objects", err)
	}
	if _, err = r.dir("refs", true); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", "refs", err)
	}
	if _, err = r.dir("refs/tags", true); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", "refs/tags", err)
	}
	if _, err = r.dir("refs/heads", true); err != nil {
		return fmt.Errorf("failed to dir path %s.\n\terror is caused by %w", "refs/heads", err)
	}

	filePath, err := r.file("description", false)
	if err != nil {
		return fmt.Errorf("failed to file path %s.\n\terror is caused by %w", "description", err)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file path %s.\n\terror is caused by %w", filePath, err)
	}
	n, err := file.WriteString("unnamed repository: edit this file 'description' to name the repository.")
	if err != nil {
		return fmt.Errorf("failed to write %s", filePath)
	}
	if n == 0 {
		return fmt.Errorf("failed to write %s. write number of bytes is 0", filePath)
	}

	filePath, err = r.file("HEAD", false)
	if err != nil {
		return fmt.Errorf("failed to file path %s.\n\terror is caused by %w", "description", err)
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file path %s.\n\terror is caused by %w", filePath, err)
	}
	n, err = file.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return fmt.Errorf("failed to write %s", filePath)
	}
	if n == 0 {
		return fmt.Errorf("failed to write %s. write number of bytes is 0", filePath)
	}

	filePath, err = r.file("config", false)
	if err != nil {
		return fmt.Errorf("failed to file path %s.\n\terror is caused by %w", "config", err)
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file path %s.\n\terror is caused by %w", filePath, err)
	}
	config := r.defaultConfig()
	n, err = file.WriteString(config)
	if err != nil {
		return fmt.Errorf("failed to write %s", filePath)
	}
	if n == 0 {
		return fmt.Errorf("failed to write %s. write number of bytes is 0", filePath)
	}

	return nil
}

func (r *repository) defaultConfig() string {
	return `[core]
		repositoryformatversion = 0
		filemode = false
		bare = false
	`
}
