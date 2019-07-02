package heartbeat

import (
	"time"
	"sync"
	"storage/rabbitmq"
	"storage/config"
	"math/rand"
	"strconv"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(config.RabbitMQServer)
	defer q.Close()

	q.Bind("apiServers")
	c:= q.Consume()

	go removeExpireDataServer()

	for msg := range c {
		dataServer , e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}

		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()

	}

}

func removeExpireDataServer() {
	for  {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s , t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers , s)
			}
		}
		mutex.Unlock()
	}
}

func GetDataServers()[]string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string , 0 )
	for s , _ := range dataServers{
		ds = append(ds , s)
	}
	return ds
}

func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return  ""
	}

	return ds[rand.Intn(n)]
}