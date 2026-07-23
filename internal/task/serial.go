/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:10
 * @Desc:
 */

package task

import (
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func GetDataFromSerial(portPath string, chanRcvSerialData chan<- string) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open(portPath, mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 1024)
	retryDelay := time.Millisecond * 100

	for {
		n, err := port.Read(buff)
		if err != nil {
			if err.Error() == "Port has been closed" {
				retryDelay *= 2
				if retryDelay > time.Second*5 {
					retryDelay = time.Second * 5
				}
				log.WithField("retry_delay", retryDelay.String()).Warn("serial port closed, reconnecting")
				time.Sleep(retryDelay)
				port, err = serial.Open(portPath, mode)
				if err != nil {
					log.WithField("err", err).Warn("serial port reopen failed")
					continue
				}
				retryDelay = time.Millisecond * 100
			} else {
				log.WithField("err", err).Warn("serial read error")
				continue
			}
		}

		if n < 1 {
			continue
		}
		chanRcvSerialData <- string(buff[:n])
		time.Sleep(time.Millisecond * 10)

	}
}
