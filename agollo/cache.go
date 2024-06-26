package agollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type namespaceCache struct {
	lock   sync.RWMutex
	caches map[string]*cache
}

func newNamespaceCahce() *namespaceCache {
	return &namespaceCache{
		caches: map[string]*cache{},
	}
}

func (n *namespaceCache) mustGetCache(namespace string) *cache {
	n.lock.RLock()
	if ret, ok := n.caches[namespace]; ok {
		n.lock.RUnlock()
		return ret
	}
	n.lock.RUnlock()

	n.lock.Lock()
	defer n.lock.Unlock()

	cache := newCache()
	n.caches[namespace] = cache
	return cache
}

func (n *namespaceCache) drain() {
	for namespace := range n.caches {
		delete(n.caches, namespace)
	}
}

func (n *namespaceCache) dump(name string) error {
	var dumps = map[string]map[string]string{}
	for namespace, cache := range n.caches {
		dumps[namespace] = cache.dump()
	}
	for namespace := range dumps {
		content := dumps[namespace]["content"]
		if dumps[namespace]["content"] != "" {
			f, err := os.OpenFile(fmt.Sprintf("%s%s", name, namespace), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			if _, err = f.Write([]byte(content)); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

func (n *namespaceCache) load(name string) error {
	n.drain()
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}

	var dumps = map[string]map[string]string{}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), name) {
			f, err := os.OpenFile(file.Name(), os.O_RDONLY, 0755)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				f.Close()
				return err
			}
			namespace := strings.Replace(file.Name(), name, "", 1)
			if _, ok := dumps[namespace]; !ok {
				dumps[namespace] = map[string]string{}
			}
			dumps[namespace]["content"] = string(content)
			f.Close()
		}
	}

	for namespace, kv := range dumps {
		tempCache := n.mustGetCache(namespace)
		for k, v := range kv {
			tempCache.set(k, v)
		}
	}

	return nil
}

// by xingdonghai
func (n *namespaceCache) decode(model interface{}, oneNamespaceMode bool, tagName string) error {
	// prepare input
	input := make(map[string]interface{})
	var ns string
	for namespace, cache := range n.caches {
		// by xingdonghai
		if ns == "" {
			ns = namespace
		}

		switch {
		case strings.Contains(namespace, ".yaml"):
			if content, ok := cache.get("content"); ok && content != "" {
				v := make(map[string]interface{})
				yaml.Unmarshal([]byte(content), &v)
				input[namespace] = v
			}
		case strings.Contains(namespace, ".json"):
			if content, ok := cache.get("content"); ok && content != "" {
				v := make(map[string]interface{})
				json.Unmarshal([]byte(content), &v)
				input[namespace] = v
			}
		default:
			if tagName == "properties" {
				if content, ok := cache.get("content"); ok && content != "" {
					v := make(map[string]interface{})
					json.Unmarshal([]byte(content), &v)
					input[namespace] = v
				}
			} else {
				input[namespace] = cache.dump()
			}
		}

	}

	// decode
	// by xingdonghai
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   model,
		TagName:  tagName,
	})
	if err != nil {
		return err
	}
	if oneNamespaceMode && len(input) == 1 {
		return decoder.Decode(input[ns])
	} else {
		return decoder.Decode(input)
	}

	// return mapstructure.Decode(input, model)
}

type cache struct {
	kv sync.Map
}

func newCache() *cache {
	return &cache{
		kv: sync.Map{},
	}
}

func (c *cache) set(key, val string) {
	c.kv.Store(key, val)
}

func (c *cache) get(key string) (string, bool) {
	if val, ok := c.kv.Load(key); ok {
		if ret, ok := val.(string); ok {
			return ret, true
		}
	}
	return "", false
}

func (c *cache) delete(key string) {
	c.kv.Delete(key)
}

func (c *cache) dump() map[string]string {
	var ret = map[string]string{}
	c.kv.Range(func(key, val interface{}) bool {
		k, _ := key.(string)
		v, _ := val.(string)
		ret[k] = v

		return true
	})
	return ret
}
