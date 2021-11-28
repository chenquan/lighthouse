/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package config

import (
	"time"
)

type Configuration interface {
	// Validate validates the configuration.
	// If returns error, the broker will not start.
	Validate() error
	// Unmarshaler defined how to unmarshal YAML into the config structure.
	//yaml.Unmarshaler
}
type Config struct {
	Mqtt Mqtt `yaml:"mqtt"`
}

func (c *Config) Validate() error {
	panic("implement me")
}

//func (c *Config) UnmarshalYAML(value *yaml.Node) error {
//	return nil
//}

type Mqtt struct {
	// SessionExpiry is the maximum session expiry interval in seconds.
	SessionExpiry time.Duration `yaml:"session_expiry"`
	// SessionExpiryCheckInterval is the interval time for session expiry checker to check whether there
	// are expired sessions.
	SessionExpiryCheckInterval time.Duration `yaml:"session_expiry_check_interval"`
	// MessageExpiry is the maximum lifetime of the message in seconds.
	// If a message in the queue is not sent in MessageExpiry time, it will be removed, which means it will not be sent to the subscriber.
	MessageExpiry time.Duration `yaml:"message_expiry"`
	// InflightExpiry is the lifetime of the "inflight" message in seconds.
	// If a "inflight" message is not acknowledged by a client in InflightExpiry time, it will be removed when the message queue is full.
	InflightExpiry time.Duration `yaml:"inflight_expiry"`
	// MaxPacketSize is the maximum packet size that the server is willing to accept from the client
	MaxPacketSize uint32 `yaml:"max_packet_size"`
	// ReceiveMax limits the number of QoS 1 and QoS 2 publications that the server is willing to process concurrently for the client.
	ReceiveMax uint16 `yaml:"server_receive_maximum"`
	// MaxKeepAlive is the maximum keep alive time in seconds allows by the server.
	// If the client requests a keepalive time bigger than MaxKeepalive,
	// the server will use MaxKeepAlive as the keepalive time.
	// In this case, if the client version is v5, the server will set MaxKeepalive into CONNACK to inform the client.
	// But if the client version is 3.x, the server has no way to inform the client that the keepalive time has been changed.
	MaxKeepAlive uint16 `yaml:"max_keepalive"`
	// TopicAliasMax indicates the highest value that the server will accept as a Topic Alias sent by the client.
	// No-op if the client version is MQTTv3.x
	TopicAliasMax uint16 `yaml:"topic_alias_maximum"`
	// SubscriptionIDAvailable indicates whether the server supports Subscription Identifiers.
	// No-op if the client version is MQTTv3.x .
	SubscriptionIDAvailable bool `yaml:"subscription_identifier_available"`
	// SharedSubAvailable indicates whether the server supports Shared Subscriptions.
	SharedSubAvailable bool `yaml:"shared_subscription_available"`
	// WildcardSubAvailable indicates whether the server supports Wildcard Subscriptions.
	WildcardAvailable bool `yaml:"wildcard_subscription_available"`
	// RetainAvailable indicates whether the server supports retained messages.
	RetainAvailable bool `yaml:"retain_available"`
	// MaxQueuedMsg is the maximum queue length of the outgoing messages.
	// If the queue is full, some message will be dropped.
	// The message dropping strategy is described in the document of the persistence/queue.Store interface.
	MaxQueueMessages int `yaml:"max_queue_messages"`
	// MaxInflight limits inflight message length of the outgoing messages.
	// Inflight message is also stored in the message queue, so it must be less than or equal to MaxQueuedMsg.
	// Inflight message is the QoS 1 or QoS 2 message that has been sent out to a client but not been acknowledged yet.
	MaxInflight uint16 `yaml:"max_inflight"`
	// MaximumQoS is the highest QOS level permitted for a Publish.
	MaximumQoS uint8 `yaml:"maximum_qos"`
	// QueueQos0Msg indicates whether to store QoS 0 message for a offline session.
	QueueQos0Msg bool `yaml:"queue_qos0_messages"`
	// DeliveryMode is the delivery mode. The possible value can be "overlap" or "onlyonce".
	// It is possible for a client’s subscriptions to overlap so that a published message might match multiple filters.
	// When set to "overlap" , the server will deliver one message for each matching subscription and respecting the subscription’s QoS in each case.
	// When set to "onlyOnce",the server will deliver the message to the client respecting the maximum QoS of all the matching subscriptions.
	DeliveryMode string `yaml:"delivery_mode"`
	// AllowZeroLenClientId indicates whether to allow a client to connect with empty client id.
	AllowZeroLenClientId bool `yaml:"allow_zero_len_client_id"`
}
