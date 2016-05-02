// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/mdns"
	"github.com/surgemq/surgemq/service"
)

func startMQTT() error {
	svr := &service.Server{
		KeepAlive:        5 * 60,
		ConnectTimeout:   2,
		AckTimeout:       20,
		TimeoutRetries:   3,
		SessionsProvider: service.DefaultSessionsProvider,
		TopicsProvider:   service.DefaultTopicsProvider,
	}
	return svr.ListenAndServe("tcp://:1883")
}

/*
func pinger() {
	p = &netx.Pinger{}
	if err := p.AddIPs(serverIPs); err != nil {
		log.Fatal(err)
	}

	cnt := 0
	tick := time.NewTicker(time.Duration(pingInterval) * time.Second)

	for {
		if cnt != 0 {
			<-tick.C
		}

		res, err := p.Start()
		if err != nil {
			log.Fatal(err)
		}

		for pr := range res {
			if !serverQuiet {
				log.Println(pr)
			}

			// Creates a new PUBLISH message with the appropriate contents for publishing
			pubmsg := message.NewPublishMessage()
			if pr.Err != nil {
				pubmsg.SetTopic([]byte(fmt.Sprintf("/ping/failure/%s", pr.Src)))
			} else {
				pubmsg.SetTopic([]byte(fmt.Sprintf("/ping/success/%s", pr.Src)))
			}
			pubmsg.SetQoS(0)

			b, err := pr.GobEncode()
			if err != nil {
				log.Printf("pinger: Error from GobEncode: %v\n", err)
				continue
			}

			pubmsg.SetPayload(b)

			// Publishes to the server
			s.Publish(pubmsg, nil)
		}

		p.Stop()
		cnt++
	}
}

func onPublish(msg *message.PublishMessage) error {
		pr := &netx.PingResult{}
		if err := pr.GobDecode(msg.Payload()); err != nil {
			log.Printf("Error decoding ping result: %v\n", err)
			return err
		}
		log.Println(pr)
	return nil
}

func client1(host string) error {
	c = &service.Client{}
	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId("hostname")
	msg.SetKeepAlive(300)
	if err := c.Connect("tcp://"+host+":1883", msg); err != nil {
		return err
	}
	submsg := message.NewSubscribeMessage()
	for _, t := range clientTopics {
		submsg.AddTopic([]byte(t), 0)
	}
	c.Subscribe(submsg, nil, onPublish)
	<-done
}
*/

func defaultHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func client(host string) error {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + host + ":1883")
	opts.SetClientID("dlibox")
	opts.SetDefaultPublishHandler(defaultHandler)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := c.Subscribe("dlibox/led", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("dlibox/led", 0, false, text)
		token.Wait()
	}
	time.Sleep(3 * time.Second)
	if token := c.Unsubscribe("dlibox/led"); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.Disconnect(250)
	return nil
}

// Respond to incoming requests and look up devices.
func initmDNS(properties []string) (*mDNS, error) {
	// Always try to start as a MQTT broker. If success (e.g. port was
	// available), list itself as server, otherwise don't.

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	// "_http._tcp."
	// http://www.dns-sd.org/servicetypes.html
	// http://www.iana.org/form/ports-services
	service, err := mdns.NewMDNSService(hostName, "dlibox", "", "", 1883, nil, properties)
	if err != nil {
		return nil, err
	}
	l, err := net.Listen("udp", ":3611")
	if err != nil {
		return nil, err
	}
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		l.Close()
		return nil, err
	}
	out := &mDNS{s: server, l: l}
	go out.listen()
	go out.lookup()
	return out, err
}

type mDNS struct {
	s *mdns.Server
	l net.Listener // Communicates over UDP. Eventually using MQTT would be a good idea.

	lock    sync.Mutex
	entries []*mdns.ServiceEntry
}

func (m *mDNS) Close() error {
	err1 := m.s.Shutdown()
	err2 := m.l.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// Packet is contextual with the From and To entries in
type Packet struct {
}

func (m *mDNS) Entries() []*mdns.ServiceEntry {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.entries
}

// IsMaster returns true if this device should be the master. The master is the
// device with the smallest serial number, as listed by the hostname.
func (m *mDNS) IsMaster() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	return true
}

func (m *mDNS) listen() {
	for {
		conn, err := m.l.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			c.Read(nil)
			c.Close()
		}(conn)
	}
}

func (m *mDNS) lookup() {
	// TODO(maruel): When another device polls for services, immediately register
	// the device too.
	entries := make(chan *mdns.ServiceEntry)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var entries []*mdns.ServiceEntry
		for _, e := range entries {
			entries = append(entries, e)
		}
		m.lock.Lock()
		m.entries = entries
		m.lock.Unlock()
	}()
	defer wg.Wait()
	defer close(entries)
	// Is 1sec enough? The ESP8266 isn't fast.
	_ = mdns.Lookup("dlibox", entries)
}
