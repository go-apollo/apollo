//Copyright (c) 2017 Phil

package apollo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		caches:         newNamespaceCache(),
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
	if err != nil {
		return fmt.Errorf("handling namespace not updated %s", err)
	}
	//don't delivery change event if namespace has any change
	if change == nil {
		return nil
	}
	//delivery change event if namespace has changed
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

// fetchAllConfig fetch from remote, if failed load from local file
func (c *Client) preload() error {
	if err := c.longPoller.preload(); err != nil {
		log.Errorf("fetchAllConfig fetch from remote: %s", err)
		return c.loadLocal(defaultDumpFile)
	}
	return nil
}

// loadLocal load caches from local file
func (c *Client) loadLocal(name string) error {
	log.Debugf("loadLocal load caches from local file,file name: %s", name)
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

func (c *Client) mustGetCache(namespace string) *cache {
	return c.caches.mustGetCache(namespace)
}

// GetStringValueWithNameSpace get value from given namespace
func (c *Client) GetStringValueWithNameSpace(namespace, key, defaultValue string) string {
	log.Debugf("GetStringValueWithNameSpace get value from given namespace,namespace: %s,key: %s, defaultValue: %s", namespace, key, defaultValue)
	cache := c.mustGetCache(namespace)
	if ret, ok := cache.get(key); ok {
		log.Debugf("GetStringValueWithNameSpace from cache result:\nret: %s \nisOk: %t", ret, ok)
		return string(ret)
	}
	return defaultValue
}

// GetIntValueWithNameSpace get int value from given namespace
func (c *Client) GetIntValueWithNameSpace(namespace, key string, defaultValue int) int {
	sValue := GetStringValueWithNameSpace(namespace, key, "")
	intValue, err := strconv.Atoi(sValue)
	if err != nil {
		log.Errorf("GetIntValue %s err: %s", key, err.Error())
		return defaultValue
	}
	return intValue
}

// GetStringValue from default namespace
func (c *Client) GetStringValue(key, defaultValue string) string {
	return c.GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// GetIntValue from default namespace
func (c *Client) GetIntValue(key string, defaultValue int) int {
	return c.GetIntValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// GetNameSpaceContent get contents of namespace
func (c *Client) GetNameSpaceContent(namespace, defaultValue string) string {
	return c.GetStringValueWithNameSpace(namespace, "content", defaultValue)
}

// ListKeys list all keys under given namespace
func (c *Client) ListKeys(namespace string) []string {
	var keys []string
	cache := c.mustGetCache(namespace)
	cache.kv.Range(func(k, _ interface{}) bool {
		str, ok := k.(string)
		if ok {
			keys = append(keys, str)
		}
		return true
	})
	return keys
}

// sync namespace config
func (c *Client) sync(namespace string) (*ChangeEvent, error) {
	releaseKey := c.getReleaseKey(namespace)
	url := configURL(c.conf, namespace, string(releaseKey))
	bts, err := c.requester.request(url)
	if err != nil || len(bts) == 0 {
		return nil, fmt.Errorf("sync namespace config error, remote error or empty congfig")
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
		cache.set(k, []byte(v))
		old, ok := kv[k]
		if !ok {
			ret.Changes[k] = makeAddChange(k, []byte(v))
			continue
		}
		if string(old) != string(v) {
			ret.Changes[k] = makeModifyChange(k, old, []byte(v))
		}
	}

	c.setReleaseKey(result.NamespaceName, []byte(result.ReleaseKey))

	// dump caches to file
	c.dump(defaultDumpFile)

	if len(ret.Changes) == 0 {
		return nil
	}

	return &ret
}

func (c *Client) getReleaseKey(namespace string) []byte {
	releaseKey, _ := c.releaseKeyRepo.get(namespace)
	return releaseKey
}

func (c *Client) setReleaseKey(namespace string, releaseKey []byte) {
	c.releaseKeyRepo.set(namespace, releaseKey)
}
