/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:10
 * @Desc:
 */

package task

import (
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"time"
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

	for {
		n, err := port.Read(buff)
		if err != nil {
			if err.Error() == "Port has been closed" {
				time.Sleep(time.Millisecond * 100)
				port, err = serial.Open(portPath, mode)
				if err != nil {
					log.WithField("err", err).Fatal()
				}
			} else {
				log.WithField("err", err).Fatal()
			}
		}

		if n < 1 {
			continue
		} else {
			chanRcvSerialData <- string(buff[:n])
		}
		time.Sleep(time.Millisecond * 10)

	}
}
