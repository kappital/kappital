/*
 * Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package watcher

import (
	"encoding/json"
	"fmt"
	"reflect"

	"k8s.io/klog/v2"
)

const (
	// OPList means query a type list
	OPList = "LIST"
	// OPDelete means delete a type data
	OPDelete = "DELETE"
	// OPCreate means insert a type data
	OPCreate = "INSERT"
	// OPUpdate means update a type data
	OPUpdate = "UPDATE"

	channalBuffer = 1024

	klogLevel = 3
)

var notifyChannel chan *NotifyInfo

// NotifyInfo stores a structure of a certain type
type NotifyInfo struct {
	Notify      Notification
	ChannelName string
}

// NotifyWatcher Event listening structure
type NotifyWatcher struct {
	NotifyMsg      chan *NotifyInfo
	listenChannels map[string]*channelConfig
	stopCh         chan struct{}
}

// Watch the events
func (n *NotifyWatcher) Watch(obj interface{}, channel string, processor Processor) error {
	if obj == nil || channel == "" || processor == nil {
		return fmt.Errorf("any of 'obj', 'channel', 'processor' cannot be emypty value")
	}

	n.listenChannels[channel] = &channelConfig{
		name:      channel,
		target:    obj,
		processor: processor,
	}
	klog.Infof("register watcher for obj %T on channel %s", obj, channel)
	return nil
}

// StartProcessor for the synchronizing, watch the events and process them
func (n *NotifyWatcher) StartProcessor() error {
	notifyChannel = make(chan *NotifyInfo, channalBuffer)
	n.NotifyMsg = notifyChannel
	n.stopCh = make(chan struct{})
	go n.dispatch()

	for _, conf := range n.listenChannels {
		objs, err := conf.processor.List()
		if err != nil {
			return err
		}
		for _, o := range objs {
			conf.processor.Process(OPList, o)
		}
	}

	return nil
}

// Start the synchronizing
func (n *NotifyWatcher) Start() error {
	return nil
}

func (n *NotifyWatcher) dispatch() {
	for {
		select {
		case m, ok := <-n.NotifyMsg:
			if !ok {
				klog.Errorf("failed to get notify msg channel")
				return
			}
			if m == nil {
				klog.Warningf("database watcher got nil notification")
				continue
			}
			klog.V(klogLevel).Infof("database watcher got notification from/%s: %s", m.ChannelName, m.Notify)
			channelConfig, ok := n.listenChannels[m.ChannelName]
			if !ok {
				klog.Warningf("watcher got notification from unknown channel %s, ignore", m.ChannelName)
				continue
			}

			objType := reflect.TypeOf(channelConfig.target)
			newObj := reflect.New(objType.Elem()).Interface()
			if err := json.Unmarshal([]byte(m.Notify.RawData), newObj); err != nil {
				klog.Errorf("unmarshal notification '%s' failed: %s, ignore it", m.Notify, err)
				continue
			}

			channelConfig.processor.Process(m.Notify.OPType, newObj)
		case _, ok := <-n.stopCh:
			if ok {
				klog.V(klogLevel).Infof("stopping database listen worker")
				return
			}
		}
	}
}

// Stop the synchronize
func (n *NotifyWatcher) Stop() {
	close(n.stopCh)
	klog.Infof("notify watcher stopped.")
}

// AddEvent add the watching event into synchronizing list
func AddEvent(obj interface{}, opType, channel string) error {
	if notifyChannel == nil {
		return nil
	}
	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("failed to marshal obj, error: %v", err)
		return err
	}

	msg := NotifyInfo{
		Notify: Notification{
			OPType:  opType,
			RawData: string(objByte),
		},
		ChannelName: channel,
	}
	klog.Infof("add event id handler %s: id %s", channel, opType)
	notifyChannel <- &msg
	return nil
}

// NewWatcher create a new NotifyWatcher
func NewWatcher() Watcher {
	return &NotifyWatcher{
		listenChannels: map[string]*channelConfig{},
	}
}
