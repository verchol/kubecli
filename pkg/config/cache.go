package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

//Cache ..
type ClusterStatus int

const (
	ClusterAvailable    ClusterStatus = 1
	ClusterNotAvailable ClusterStatus = 2
	ClusterNotTested    ClusterStatus = 0
)

type KubeContext struct {
	Name           string
	Namespace      string
	Status         ClusterStatus
	AuthProvider   string
	LastUpdated    string
	CurrentContext bool
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
const ResetCache = "10" // in minutes

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
	_, err := os.Stat(c.cachePath)
	log.Printf("cachefile %v err %v", c.cachePath, err)

	if err == nil || !os.IsNotExist(err) {
		return err
	}

	log.Printf("\ncreating cache file %v\n\n", LocalCacheFile)
	err = os.MkdirAll(path.Dir(LocalCacheFile), os.ModePerm)
	if err != nil {
		return err
	}
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
	log.Printf("flashing cache\n")
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
