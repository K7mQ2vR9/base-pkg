package rabbitmq

import (
	"context"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// brokerChannel 使用 amqp091-go 的原生方法，供本包统一管理发布确认、退回及订阅恢复。
// 不启用驱动自动恢复，避免连接、确认序号及消费者被两套恢复机制同时管理。
type brokerChannel interface {
	QueueDeclare(string, bool, bool, bool, bool, amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(string, string, bool, bool, bool, bool, amqp.Table) error
	QueueBind(string, string, string, bool, amqp.Table) error
	Confirm(bool) error
	NotifyPublish(chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(chan amqp.Return) chan amqp.Return
	NotifyClose(chan *amqp.Error) chan *amqp.Error
	NotifyCancel(chan string) chan string
	Qos(int, int, bool) error
	Consume(string, string, bool, bool, bool, bool, amqp.Table) (<-chan amqp.Delivery, error)
	Publish(string, string, bool, bool, amqp.Publishing) error
}

type brokerConnection interface {
	Channel() (brokerChannel, error)
	NotifyClose(chan *amqp.Error) chan *amqp.Error
}

type amqpConnection struct{ *amqp.Connection }

func (c amqpConnection) Channel() (brokerChannel, error) { return c.Connection.Channel() }

type dialFunc func(context.Context) (brokerConnection, io.Closer, error)

// deadlineConn 为 AMQP 写操作设置时限，关闭底层连接可中断握手及阻塞的 RPC。
type deadlineConn struct {
	net.Conn
	timeout time.Duration
}

func (c *deadlineConn) Write(data []byte) (int, error) {
	if err := c.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}
	return c.Conn.Write(data)
}

// amqpDialer 使用 amqp091-go 建立受启动上下文约束的连接，恢复和关闭均由本包会话管理。
// Heartbeat 统一取已归一化的配置；通过 URI 显式传值，避免驱动将 0 再次替换为默认 10s。
func amqpDialer(config settings) dialFunc {
	return func(ctx context.Context) (brokerConnection, io.Closer, error) {
		connectionURL, err := url.Parse(config.URL)
		if err != nil {
			return nil, nil, &operationError{"parse connection URL", err}
		}
		query := connectionURL.Query()
		query.Set("heartbeat", strconv.FormatInt(int64(config.heartbeat/time.Second), 10))
		connectionURL.RawQuery = query.Encode()
		var raw net.Conn
		var stop func() bool
		conn, err := amqp.DialConfig(connectionURL.String(), amqp.Config{
			Heartbeat: config.heartbeat,
			// nil 明确关闭新驱动的自动重连与拓扑恢复，由本包重建完整会话并重新注册通知。
			Recovery: nil,
			Dial: func(network, addr string) (net.Conn, error) {
				connection, err := (&net.Dialer{Timeout: config.DialTimeout}).DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				raw = &deadlineConn{Conn: connection, timeout: config.ConfirmTimeout}
				stop = context.AfterFunc(ctx, func() { _ = raw.Close() })
				deadline, _ := ctx.Deadline()
				if err := raw.SetDeadline(deadline); err != nil {
					_ = raw.Close()
					return nil, err
				}
				return raw, nil
			},
		})
		if stop != nil {
			stop()
		}
		if err != nil {
			if raw != nil {
				_ = raw.Close()
			}
			return nil, nil, &operationError{"connect", err}
		}
		return amqpConnection{conn}, raw, nil
	}
}
