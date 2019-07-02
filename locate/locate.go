package locate

import (
	"storage/rabbitmq"
	"storage/config"
	"time"
	"os"
	"strconv"
)

func StartHeartbeat() {
	q := rabbitmq.New(config.RabbitMQServer)
	defer q.Close()
	for {
		q.Publish("apiServers" , config.RabbitMQServer)
		time.Sleep(5 * time.Second)
	}
}

func Locate(name string) bool {
	_ , err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(config.RabbitMQServer)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()

	for msg := range c{
		object , e := strconv.Unquote(string(msg.Body))
		if e != nil{
			panic(e)
		}
		if Locate(config.STORAGE_ROOT + "/objects/" + object) {
			q.Send(msg.ReplyTo , config.LISTEN_ADDRESS)
		}
	}
}