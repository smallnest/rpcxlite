package client

import (
	"context"
	"net/url"
	"strconv"

	"github.com/valyala/fastrand"
)

type SelectFunc func(ctx context.Context, servicePath, serviceMethod string, args interface{}) string

// Selector defines selector that selects one service from candidates.
type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // SelectFunc
	UpdateServer(servers map[string]string)
}

func newSelector(selectMode SelectMode, servers map[string]string) Selector {
	switch selectMode {
	case RandomSelect:
		return newRandomSelector(servers)
	case RoundRobin:
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case SelectByUser:
		return nil
	default:
		return newRandomSelector(servers)
	}
}

// randomSelector selects randomly.
type randomSelector struct {
	servers []string
}

func newRandomSelector(servers map[string]string) Selector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	return &randomSelector{servers: ss}
}

func (s randomSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := fastrand.Uint32n(uint32(len(ss)))
	return ss[i]
}

func (s *randomSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}

// roundRobinSelector selects servers with roundrobin.
type roundRobinSelector struct {
	servers []string
	i       int
}

func newRoundRobinSelector(servers map[string]string) Selector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	return &roundRobinSelector{servers: ss}
}

func (s *roundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := s.i
	i = i % len(ss)
	s.i = i + 1

	return ss[i]
}

func (s *roundRobinSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}

// weightedRoundRobinSelector selects servers with weighted.
type weightedRoundRobinSelector struct {
	servers []*Weighted
}

func newWeightedRoundRobinSelector(servers map[string]string) Selector {
	ss := createWeighted(servers)
	return &weightedRoundRobinSelector{servers: ss}
}

func (s *weightedRoundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	w := nextWeighted(ss)
	if w == nil {
		return ""
	}
	return w.Server
}

func (s *weightedRoundRobinSelector) UpdateServer(servers map[string]string) {
	ss := createWeighted(servers)
	s.servers = ss
}

func createWeighted(servers map[string]string) []*Weighted {
	ss := make([]*Weighted, 0, len(servers))
	for k, metadata := range servers {
		w := &Weighted{Server: k, Weight: 1, EffectiveWeight: 1}

		if v, err := url.ParseQuery(metadata); err == nil {
			ww := v.Get("weight")
			if ww != "" {
				if weight, err := strconv.Atoi(ww); err == nil {
					w.Weight = weight
					w.EffectiveWeight = weight
				}
			}
		}

		ss = append(ss, w)
	}

	return ss
}
