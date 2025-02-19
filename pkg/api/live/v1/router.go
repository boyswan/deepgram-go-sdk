// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the live API for Deepgram
package live

import (
	"encoding/json"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// MessageRouter is helper struct that routes events
type MessageRouter struct {
	callback interfaces.LiveMessageCallback
}

// NewWithDefault creates a default MessageRouter
func NewWithDefault() *MessageRouter {
	return New(NewDefaultCallbackHandler())
}

// New creates a MessageRouter with user defined callback
func New(callback interfaces.LiveMessageCallback) *MessageRouter {
	return &MessageRouter{
		callback: callback,
	}
}

// Message handles platform messages
func (r *MessageRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("router.Message ENTER\n")

	// what is the high level message here?
	var mt MessageType
	err := json.Unmarshal(byMsg, &mt)
	if err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("router.Message LEAVE\n")
		return err
	}

	klog.V(6).Infof("router.Message LEAVE\n")

	switch mt.Type {
	case interfaces.TypeErrorResponse:
		return r.ErrorResponse(byMsg)
	case interfaces.TypeMessageResponse:
		return r.MessageResponse(byMsg)
	case interfaces.TypeMetadataResponse:
		return r.MetadataResponse(byMsg)
	case interfaces.TypeUtteranceEndResponse:
		return r.UtteranceEndResponse(byMsg)
	default:
		return r.UnhandledMessage(byMsg)
	}
}

// MessageResponse handles the MessageResponse message
func (r *MessageRouter) MessageResponse(byMsg []byte) error {
	klog.V(6).Infof("router.MessageResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "MessageResponse", byMsg)

	var mr interfaces.MessageResponse
	err := json.Unmarshal(byMsg, &mr)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.MessageResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Message(&mr)
		if err != nil {
			klog.V(1).Infof("callback.MessageResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.MessageResponse succeeded\n")
		}
		klog.V(6).Infof("router.MessageResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.MessageResponse LEAVE\n")

	return ErrUserCallbackNotDefined
}

func (r *MessageRouter) MetadataResponse(byMsg []byte) error {
	klog.V(6).Infof("router.MetadataResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "MetadataResponse", byMsg)

	var md interfaces.MetadataResponse
	err := json.Unmarshal(byMsg, &md)
	if err != nil {
		klog.V(1).Infof("MetadataResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.MetadataResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Metadata(&md)
		if err != nil {
			klog.V(1).Infof("callback.MetadataResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.MetadataResponse succeeded\n")
		}
		klog.V(6).Infof("router.MetadataResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.MetadataResponse ENTER\n")

	return nil
}

func (r *MessageRouter) UtteranceEndResponse(byMsg []byte) error {
	klog.V(6).Infof("router.UtteranceEnd ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "UtteranceEnd", byMsg)

	var ur interfaces.UtteranceEndResponse
	err := json.Unmarshal(byMsg, &ur)
	if err != nil {
		klog.V(1).Infof("UtteranceEnd json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.UtteranceEnd LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.UtteranceEnd(&ur)
		if err != nil {
			klog.V(1).Infof("callback.UtteranceEnd failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.UtteranceEnd succeeded\n")
		}
		klog.V(6).Infof("router.UtteranceEnd LEAVE\n")
		return err
	}

	return nil
}

func (r *MessageRouter) ErrorResponse(byMsg []byte) error {
	klog.V(6).Infof("router.ErrorResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "ErrorResponse", byMsg)

	var er interfaces.ErrorResponse
	err := json.Unmarshal(byMsg, &er)
	if err != nil {
		klog.V(1).Infof("ErrorResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.ErrorResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Error(&er)
		if err != nil {
			klog.V(1).Infof("callback.ErrorResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.ErrorResponse succeeded\n")
		}
		klog.V(6).Infof("router.ErrorResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.ErrorResponse ENTER\n")

	return nil
}

// UnhandledMessage handles the UnhandledMessage message
func (r *MessageRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("router.UnhandledMessage ENTER\n")

	// trace debugging
	r.printDebugMessages(2, "UnhandledMessage", byMsg)

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

func (r *MessageRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
	}

	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJson)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
