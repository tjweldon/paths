package fullpath

import (
	"encoding/json"
	"os"
	"strings"
)

type Paths struct {
	Paths []string
}

func (p *Paths) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Paths)
}

func (p *Paths) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &p.Paths)
	return err
}

func (p *Paths) ReadConfig(configPath string) (*Paths, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(content, p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (p *Paths) ReadEnv() *Paths {
	paths := os.Getenv(pathKey)
	p.Paths = strings.Split(paths, ":")
	return p
}

func (p *Paths) Deduplicate() *Paths {
	duped := p.Paths

	p.Paths = make([]string, 0, len(duped))
	for _, path := range duped {
		if stringInSlice(path, p.Paths) {
			continue
		}
		p.Paths = append(p.Paths, path)
	}
	return p
}

func (p *Paths) Insert(path string, index int) {
	if index < 0 {
		index = 0
	}
	if index >= len(p.Paths) {
		index = len(p.Paths)
	}

	segOne := p.Paths[:index]
	segTwo := append([]string{path}, p.Paths[index:]...)
	p.Paths = append(segOne, segTwo...)
}

func (p *Paths) Remove(index int) {
	if index < 0 {
		index = 0
	}
	if index >= len(p.Paths) {
		index = len(p.Paths) - 1
	}

	segOne := p.Paths[:index]
	segTwo := append([]string{}, p.Paths[index+1:]...)
	p.Paths = append(segOne, segTwo...)
}

func (p *Paths) Move(src, dst int) {
	if src >= len(p.Paths) || src < 0 || src == dst {
		return
	}

	path := p.Paths[src]
	p.Remove(src)
	p.Insert(path, dst)
}

func (p *Paths) Replace(target int, path string) {
	if target < 0 || target >= len(p.Paths) {
		return
	}

	p.Paths[target] = path
}

func (p *Paths) Swap(src, dst int) {
	if src == dst {
		return
	}

	srcPath, dstPath := p.Paths[src], p.Paths[dst]
	p.Replace(src, dstPath)
	p.Replace(dst, srcPath)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

const pathKey = "PATH"
