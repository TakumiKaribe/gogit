package wyago

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

func new(path string, force bool) (*repository, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to new repository. error is caused by %w", err)
	}
	if !force && !fileInfo.IsDir() {
		return nil, fmt.Errorf("failed to new repository. %s is not a git repository", path)
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

	_, err = r.file("config", false)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to new repository. error is caused by %w", err)
	}
	if err != nil && os.IsNotExist(err) && !force {
		return nil, fmt.Errorf("failed to new repository. config file missing")
	}

	if !force {
		err := r.viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to new repository. error is caused by %w", err)
		}
		version, ok := r.viper.GetStringMap("core")["repositoryformatversion"]
		if !ok {
			return nil, fmt.Errorf("failed to new repository. repositoryformatversion is not found")
		}
		_version, ok := version.(int)
		if !ok {
			return nil, fmt.Errorf("failed to new repository. repositoryformatversion is not integer")
		}
		if _version != 0 {
			return nil, fmt.Errorf("failed to new repository. unsupported repositoryformatversion %v", _version)
		}
	}

	return r, nil
}

func (r *repository) path(path string) string {
	return r.gitDir + path
}

func (r *repository) file(path string, mkdir bool) (string, error) {
	path, err := r.dir(path, mkdir)
	if err != nil {
		return "", err
	}
	return r.path(path), nil
}

func (r *repository) dir(path string, mkdir bool) (string, error) {
	path = r.path(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if fileInfo.IsDir() {
				return path, nil
			}
			return "", fmt.Errorf("failed dir. %s is not a directory", path)
		}
		return "", err
	}

	if mkdir {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	return "", nil
}

func (r *repository) create(path string) error {
	fileInfo, err := os.Stat(r.worktree)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(r.worktree, os.ModeDir)
		if err != nil {
			return fmt.Errorf("failed to create repository. error is caused by %w", err)
		}
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("failed to create repository. %s is not a directory", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if len(files) != 0 {
		return fmt.Errorf("failed to create repository. %s is not empty", path)
	}

	if _, err = r.dir("branches", true); err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if _, err = r.dir("objects", true); err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if _, err = r.dir("refs", true); err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if _, err = r.dir("refs/tags", true); err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if _, err = r.dir("refs/heads", true); err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}

	filePath, err := r.file("description", false)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	n, err := file.WriteString("unnamed repository: edit this file 'description' to name the repository.")
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if n == 0 {
		return fmt.Errorf("failed to create repository. failed to write %s", filePath)
	}

	filePath, err = r.file("HEAD", false)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	n, err = file.WriteString("unnamed repository: edit this file 'description' to name the repository.")
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if n == 0 {
		return fmt.Errorf("failed to create repository. failed to write %s", filePath)
	}

	filePath, err = r.file("config", false)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	config := r.defaultConfig()
	n, err = file.WriteString(config)
	if err != nil {
		return fmt.Errorf("failed to create repository. error is caused by %w", err)
	}
	if n == 0 {
		return fmt.Errorf("failed to create repository. failed to write %s", filePath)
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
