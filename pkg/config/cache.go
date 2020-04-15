package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//Cache ..

type KubeContext struct {
	Name         string
	Namespace    string
	Status       bool
	AuthProvider string
	LastUpdated  string
}
type Cache interface {
	Create() error
	//Empty() error
	//Delete() error
	AddEntry(entry string, c *KubeContext)
	GetEntry(entry string) (*KubeContext, error)
	Reset() bool
}

const LocalCacheFileName = ".status-test-cache"

var LocalCacheFile string

const LocalCacheSize = 100

type LocalCache struct {
	cachePath string
	cache     map[string]*KubeContext
}

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	LocalCacheFile = path.Join(homedir, ".kubecli", LocalCacheFileName)
}

//NewLocalCache ...
func NewLocalCache(cacheOpt ...string) (*LocalCache, error) {
	var cachePath string
	if len(cacheOpt) == 0 || cacheOpt[0] == string("") {
		cachePath = LocalCacheFile
	} else {
		cachePath = cacheOpt[0]
	}

	c := &LocalCache{cachePath: cachePath}
	_, err := c.loadCache()
	if err != nil {
		return c, err
	}
	//c.cache = cache
	return c, nil
}

//Create ..
func (c *LocalCache) Create() error {
	fileInfo, err := os.Stat(c.cachePath)
	fmt.Printf("cachefile %v err %v", c.cachePath, err)

	if fileInfo.Size() != 0 {
		return nil
	}

	fmt.Printf("creating cache file %v", LocalCacheFile)
	_, err = os.Create(LocalCacheFile)

	return err
}
func (c *LocalCache) loadCache() (map[string]*KubeContext, error) {
	err := c.Create()
	if err != nil {
		panic(err)
	}
	cache := make(map[string]*KubeContext, LocalCacheSize)
	bytes, err := ioutil.ReadFile(LocalCacheFile)
	json.Unmarshal(bytes, &cache)
	if err != nil {
		return nil, err
	}
	c.cache = cache
	return cache, nil
}

func (c *LocalCache) Flash() (*LocalCache, error) {

	bytes, err := json.Marshal(c.cache)

	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(c.cachePath, bytes, 0644)

	return c, err
}
func (c *LocalCache) AddEntry(entry string, ctx *KubeContext) *LocalCache {

	c.cache[entry] = ctx

	return c
}

func (c *LocalCache) GetEntry(entry string) (*KubeContext, error) {
	cacheObj, err := c.loadCache()
	if err != nil {
		return nil, err
	}
	kubeContext, ok := cacheObj[entry]
	if !ok {
		err = errors.New("missing key")
	}

	return kubeContext, err
}

//Reset ...
func (c *LocalCache) Reset() bool {
	c.cache = make(map[string]*KubeContext, LocalCacheSize)
	_, err := c.Flash()

	return err != nil
}
