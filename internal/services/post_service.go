package services

import (
	"YadnexTelegram/internal/configs"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

type PostConsumer struct {
	imapClient *client.Client
	config     configs.PostConfig
}

//Connecting with mail.
func (y *PostConsumer) Connect() error {
	var err error
	y.imapClient, err = client.DialTLS(y.config.Addr, nil)
	if err != nil {
		return err
	}
	if err := y.imapClient.Login(y.config.Login, y.config.Pass); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	return nil
}

//Listen for INBOX updates. Need to reconnect after idle stop.
func (y *PostConsumer) Listen(handler MessageHandler) error {
	for {
		err := y.Connect()
		if err != nil {
			return err
		}
		err = y.listen(handler)
		if err != nil {
			return err
		}
	}
}

func (y *PostConsumer) listen(handler MessageHandler) error {
	mbox, err := y.imapClient.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}

	updates := make(chan client.Update)
	y.imapClient.Updates = updates

	// Start idling
	stop := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- y.imapClient.Idle(stop, nil)
	}()

	stopped := false

	// Listen for updates
	for {
		select {
		case update := <-updates:
			switch update.(type) {
			case *client.MailboxUpdate:
				if !stopped {
					close(stop)
					fmt.Println("INBOX Mail box was updated")
					stopped = true
				}
			}
		case err := <-done:
			if err != nil {
				panic(err)
			}
			m := y.getMessage(mbox)
			fmt.Println("Message was received")
			var addrs string
			for _, from := range m.Envelope.From {
				addrs += fmt.Sprintf("%s@%s ", from.MailboxName, from.HostName)
			}
			handler.Handle(PostMessage{
				MailAddr: addrs,
				Subject:  m.Envelope.Subject,
			})
			return err
		}
	}
}

func (y *PostConsumer) getMessage(mbox *imap.MailboxStatus) *imap.Message {
	to := mbox.Messages
	seqset := new(imap.SeqSet)
	seqset.AddRange(to-1, to)

	messages := make(chan *imap.Message, 1)

	done := make(chan error, 1)
	go func() {
		done <- y.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var res *imap.Message
	for msg := range messages {
		res = msg
	}

	if err := <-done; err != nil {
		panic(err)
	}
	return res
}
