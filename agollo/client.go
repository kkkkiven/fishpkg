package agollo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client for apollo
type Client struct {
	conf *Conf

	updateChan chan *ChangeEvent

	caches         *namespaceCache
	releaseKeyRepo *cache

	longPoller poller
	requester  requester

	ctx    context.Context
	cancel context.CancelFunc
}

// result of query config
type result struct {
	// AppID          string            `json:"appId"`
	// Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

// NewClient create client from conf
func NewClient(conf *Conf) *Client {
	client := &Client{
		conf:           conf,
		caches:         newNamespaceCahce(),
		releaseKeyRepo: newCache(),

		requester: newHTTPRequester(&http.Client{Timeout: queryTimeout}),
	}

	client.longPoller = newLongPoller(conf, longPoolInterval, client.handleNamespaceUpdate)
	client.ctx, client.cancel = context.WithCancel(context.Background())
	return client
}

// Start sync config
func (c *Client) Start() error {

	// preload all config to local first
	if err := c.preload(); err != nil {
		return err
	}

	// start fetch update
	go c.longPoller.start()

	return nil
}

// handleNamespaceUpdate sync config for namespace, delivery changes to subscriber
func (c *Client) handleNamespaceUpdate(namespace string) error {
	change, err := c.sync(namespace)
	if err != nil || change == nil {
		return err
	}

	c.deliveryChangeEvent(change)
	return nil
}

// Stop sync config
func (c *Client) Stop() error {
	c.longPoller.stop()
	c.cancel()
	// close(c.updateChan)
	c.updateChan = nil
	return nil
}

// fetchAllCinfig fetch from remote, if failed load from local file
func (c *Client) preload() error {
	if debug {
		return c.loadLocal(c.getDumpFileName())
	}
	if err := c.longPoller.preload(); err != nil {
		return err
	}
	return nil
}

// loadLocal load caches from local file
func (c *Client) loadLocal(name string) error {
	return c.caches.load(name)
}

// dump caches to file
func (c *Client) dump(name string) error {
	return c.caches.dump(name)
}

// WatchUpdate get all updates
func (c *Client) WatchUpdate() <-chan *ChangeEvent {
	if c.updateChan == nil {
		c.updateChan = make(chan *ChangeEvent)
	}
	return c.updateChan
}

// OnConfigChange when config changed, run would be called
func (c *Client) OnConfigChange(run func(*ChangeEvent)) {
	go func() {
		events := c.WatchUpdate()
		for {
			select {
			case e := <-events:
				run(e)
			}
		}
	}()
}

func (c *Client) mustGetCache(namespace string) *cache {
	return c.caches.mustGetCache(namespace)
}

// GetStringValueWithNameSpace get value from given namespace
func (c *Client) GetStringValueWithNameSpace(namespace, key, defaultValue string) string {
	cache := c.mustGetCache(namespace)
	if ret, ok := cache.get(key); ok && ret != "" {
		return ret
	}
	return defaultValue
}

// GetStringValue from default namespace
func (c *Client) GetStringValue(key, defaultValue string) string {
	return c.GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// GetNameSpaceContent get contents of namespace
func (c *Client) GetNameSpaceContent(namespace, defaultValue string) string {
	return c.GetStringValueWithNameSpace(namespace, "content", defaultValue)
}

// GetAllKeys return all config keys in given namespace
func (c *Client) GetAllKeys(namespace string) []string {
	var keys []string
	cache := c.mustGetCache(namespace)
	cache.kv.Range(func(key, value interface{}) bool {
		str, ok := key.(string)
		if ok {
			keys = append(keys, str)
		}
		return true
	})
	return keys
}

// Unmarshal unmarshals the config into a struct
func (c *Client) Unmarshal(model interface{}) error {
	return c.caches.decode(model, c.conf.OneNamespaceMode, c.conf.TagName)
}

// sync namespace config
func (c *Client) sync(namespace string) (*ChangeEvent, error) {
	releaseKey := c.getReleaseKey(namespace)
	url := configURL(c.conf, namespace, releaseKey)
	bts, err := c.requester.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var result result
	if err := json.Unmarshal(bts, &result); err != nil {
		return nil, err
	}

	return c.handleResult(&result), nil
}

// deliveryChangeEvent push change to subscriber
func (c *Client) deliveryChangeEvent(change *ChangeEvent) {
	if c.updateChan == nil {
		return
	}
	select {
	case <-c.ctx.Done():
	case c.updateChan <- change:
	}
}

// handleResult generate changes from query result, and update local cache
func (c *Client) handleResult(result *result) *ChangeEvent {
	var ret = ChangeEvent{
		Namespace: result.NamespaceName,
		Changes:   map[string]*Change{},
	}

	cache := c.mustGetCache(result.NamespaceName)
	kv := cache.dump()

	for k, v := range kv {
		if _, ok := result.Configurations[k]; !ok {
			cache.delete(k)
			ret.Changes[k] = makeDeleteChange(k, v)
		}
	}

	for k, v := range result.Configurations {
		cache.set(k, v)
		old, ok := kv[k]
		if !ok {
			ret.Changes[k] = makeAddChange(k, v)
			continue
		}
		if old != v {
			ret.Changes[k] = makeModifyChange(k, old, v)
		}
	}

	if c.conf.TagName == "properties" {
		bts, err := json.Marshal(result.Configurations)
		if err == nil {
			cache.set("content", string(bts))
		}
	}

	c.setReleaseKey(result.NamespaceName, result.ReleaseKey)

	// dump caches to file
	if c.conf.EnableLocalCache {
		c.dump(c.getDumpFileName())
	}

	if len(ret.Changes) == 0 {
		return nil
	}

	return &ret
}

func (c *Client) getDumpFileName() string {
	return fmt.Sprintf("%s_%s_", c.conf.AppID, c.conf.Cluster)
}

func (c *Client) getReleaseKey(namespace string) string {
	releaseKey, _ := c.releaseKeyRepo.get(namespace)
	return releaseKey
}

func (c *Client) setReleaseKey(namespace, releaseKey string) {
	c.releaseKeyRepo.set(namespace, releaseKey)
}
