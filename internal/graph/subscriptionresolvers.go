package graph

import (
	"context"
	"postsandcomments/internal/graph/model"
	"sync"
)

type CommentEvent struct {
	PostID  string
	Comment *model.Comment
}

type SubscriptionManager struct {
	mutex       sync.RWMutex
	subscribers map[string][]chan *model.Comment
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[string][]chan *model.Comment),
		mutex:       sync.RWMutex{},
	}
}

func (m *SubscriptionManager) Subscribe(postID string) <-chan *model.Comment {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ch := make(chan *model.Comment, 1)
	m.subscribers[postID] = append(m.subscribers[postID], ch)
	return ch
}

func (m *SubscriptionManager) Unsubscribe(postID string, ch <-chan *model.Comment) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	subscribers := m.subscribers[postID]
	for i := range subscribers {
		if subscribers[i] == ch {
			subscribers = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
	m.subscribers[postID] = subscribers
}

func (m *SubscriptionManager) Publish(event *CommentEvent) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	subscribers := m.subscribers[event.PostID]
	for _, ch := range subscribers {
		ch <- event.Comment
	}
}

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := r.SubscriptionManager.Subscribe(postID)

	go func() {
		<-ctx.Done()
		r.SubscriptionManager.Unsubscribe(postID, ch)
	}()

	r.Logger.Infof("added client to subscribers for post with id = %s", postID)
	return ch, nil
}
