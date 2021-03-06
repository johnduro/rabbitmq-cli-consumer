package config

import (
	"fmt"
	"net/url"
	"path/filepath"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	RabbitMq struct {
		AmqpUrl     string
		Host        string
		Username    string
		Password    string
		Port        string
		Vhost       string
		Queue       string
		Compression bool
		Onfailure   int
	}
	Prefetch struct {
		Count  int
		Global bool
	}
	QueueSettings struct {
		Routingkey           string
		MessageTTL           int
		DeadLetterExchange   string
		DeadLetterRoutingKey string
	}
	Exchange struct {
		Name       string
		Autodelete bool
		Type       string
		Durable    bool
	}
	Logs struct {
		Error string
		Info  string
	}
}

func (c *Config) AmqpUrl() string {
	if len(c.RabbitMq.AmqpUrl) > 0 {
		return c.RabbitMq.AmqpUrl
	}

	host := c.RabbitMq.Host
	if len(c.RabbitMq.Port) > 0 {
		host = fmt.Sprintf("%s:%s", host, c.RabbitMq.Port)
	}

	uri := url.URL{
		Scheme: "amqp",
		Host:   host,
		Path:   c.RabbitMq.Vhost,
	}

	if len(c.RabbitMq.Username) > 0 {
		if len(c.RabbitMq.Password) > 0 {
			uri.User = url.UserPassword(c.RabbitMq.Username, c.RabbitMq.Password)
		} else {
			uri.User = url.User(c.RabbitMq.Username)
		}
	}

	c.RabbitMq.AmqpUrl = uri.String()

	return c.RabbitMq.AmqpUrl
}

// HasExchange checks if an exchange is configured.
func (c Config) HasExchange() bool {
	return c.Exchange.Name != ""
}

// ExchangeName returns the name of the configured exchange.
func (c Config) ExchangeName() string {
	return transformToStringValue(c.Exchange.Name)
}

// ExchangeType checks the configuration and returns the appropriate exchange type.
func (c Config) ExchangeType() string {
	// Check for missing exchange settings to preserve BC
	if "" == c.Exchange.Name && "" == c.Exchange.Type && !c.Exchange.Durable && !c.Exchange.Autodelete {
		return "direct"
	}

	return c.Exchange.Type
}

// PrefetchCount returns the configured prefetch count of the QoS settings.
func (c Config) PrefetchCount() int {
	// Attempt to preserve BC here
	if c.Prefetch.Count == 0 {
		return 3
	}

	return c.Prefetch.Count
}

// HasMessageTTL checks if a message TTL is configured.
func (c Config) HasMessageTTL() bool {
	return c.QueueSettings.MessageTTL > 0
}

// MessageTTL returns the configured message TTL.
func (c Config) MessageTTL() int32 {
	return int32(c.QueueSettings.MessageTTL)
}

// RoutingKey returns the configured key for message routing.
func (c Config) RoutingKey() string {
	return transformToStringValue(c.QueueSettings.Routingkey)
}

// HasDeadLetterExchange checks if a dead letter exchange is configured.
func (c Config) HasDeadLetterExchange() bool {
	return c.QueueSettings.DeadLetterExchange != ""
}

// DeadLetterExchange returns the configured dead letter exchange name.
func (c Config) DeadLetterExchange() string {
	return transformToStringValue(c.QueueSettings.DeadLetterExchange)
}

// HasDeadLetterRouting checks if a dead letter routing key is configured.
func (c Config) HasDeadLetterRouting() bool {
	return c.QueueSettings.DeadLetterRoutingKey != ""
}

// DeadLetterRoutingKey returns the configured key for the dead letter routing.
func (c Config) DeadLetterRoutingKey() string {
	return transformToStringValue(c.QueueSettings.DeadLetterRoutingKey)
}

func LoadAndParse(location string) (*Config, error) {
	if !filepath.IsAbs(location) {
		location, err := filepath.Abs(location)

		if err != nil {
			return nil, err
		}

		location = location
	}

	cfg := Config{}
	if err := gcfg.ReadFileInto(&cfg, location); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func CreateFromString(data string) (*Config, error) {
	cfg := &Config{}
	if err := gcfg.ReadStringInto(cfg, data); err != nil {
		return nil, err
	}

	return cfg, nil
}

func transformToStringValue(val string) string {
	if val == "<empty>" {
		return ""
	}

	return val
}
