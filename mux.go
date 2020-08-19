package supersense

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// Mux is a necessary struct to join different sources
type Mux struct {
	pipelines []chan *Event
	sources []Source
	mu *sync.Mutex
}

// NewMux returns a new mux
func NewMux(sources ...Source) (*Mux, error) {
	channels := make([]chan *Event, 0)
	m := &Mux{pipelines: channels, sources: sources, mu: &sync.Mutex{}}
	return m, nil
}

// RunAllSources run all the sources at the same time
func (m *Mux) RunAllSources(ctx context.Context) error {
	for _, s := range m.sources {
		go func(m *Mux, s Source) {
			for event := range s.Pipeline(ctx) {
				for _, pipe := range m.pipelines {
					pipe <- &event
				}
			}
		}(m, s)
	}
	for _, s := range m.sources {
		if err := s.Run(ctx); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Register attach a new channel to the pipes list.
func (m *Mux) Register(pipeline chan *Event, done <- chan struct{}) {
	m.mu.Lock()
	m.pipelines = append(m.pipelines, pipeline)
	m.mu.Unlock()
	<- done
	for i, p := range m.pipelines {
		if p == pipeline {
			m.mu.Lock()
			close(p)
			m.pipelines = append(m.pipelines[:i], m.pipelines[i+1:]...)
			m.mu.Unlock()
		}
	}
}