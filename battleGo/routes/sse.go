package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Event struct {
	ID    []byte
	Event []byte
	Data  []byte
}

type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan Event

	// New client connections
	newClients chan chan Event

	// Closed client connections
	closingClients chan chan Event

	// Client connections registry
	clients map[chan Event]bool
}

// Broker factory
func NewServer() (broker *Broker) {
	// Instantiate a EventBroker
	broker = &Broker{
		Notifier:       make(chan Event, 1),
		newClients:     make(chan chan Event),
		closingClients: make(chan chan Event),
		clients:        make(map[chan Event]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set the headers related to event streaming.
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan Event)

	// Signal the EventBroker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	ctx := req.Context()

	go func() {
		<-ctx.Done()
		broker.closingClients <- messageChan
	}()

	// block waiting for messages broadcast on this connection's messageChan
	for {
		var sb strings.Builder
		msg := <-messageChan
		if len(msg.Event) > 0 {
			sb.WriteString("event: ")
			sb.Write(msg.Event)
			sb.WriteRune('\n')
		}
		if len(msg.Data) > 0 {
			sb.WriteString("data: ")
			sb.Write(msg.Data)
			sb.WriteRune('\n')
		}
		if len(msg.ID) > 0 {
			sb.WriteString("id: ")
			sb.Write(msg.ID)
			sb.WriteRune('\n')
		}
		sb.WriteRune('\n')

		fmt.Printf("%s", sb.String())

		// Write to the ResponseWriter
		// Server Sent Events compatible
		_, _ = fmt.Fprintf(rw, "%s", sb.String())

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

// Listen on different channels and act accordingly
func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))

		case s := <-broker.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))

		case event := <-broker.Notifier:
			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}
