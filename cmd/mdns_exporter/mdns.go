package main

import (
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"codeberg.org/cheesestraws/mdns"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/logutil"
)

type mdnsServiceTracker struct {
	sync.RWMutex
	m map[string]time.Time
}

func (m *mdnsServiceTracker) contains(s string) bool {
	if m.m == nil {
		return false
	}

	m.RLock()
	defer m.RUnlock()

	_, ok := m.m[s]
	return ok
}

func (m *mdnsServiceTracker) observe(s string) bool {
	m.Lock()
	defer m.Unlock()

	if m.m == nil {
		m.m = make(map[string]time.Time)
	}

	_, replaced := m.m[s]
	m.m[s] = time.Now()

	return replaced
}

func (m *mdnsServiceTracker) removeStaleJunk(timeout time.Duration) []string {
	var removals []string

	m.Lock()
	defer m.Unlock()

	for k, v := range m.m {
		if time.Since(v) > timeout {
			removals = append(removals, k)
			delete(m.m, k)
		}
	}

	return removals
}

type mdnsWatcherStats struct {
	runningServicePollers atomic.Int64
	servicePollCount      atomic.Int64
	serviceReplyCount     atomic.Int64
	serviceCleanupCount   atomic.Int64

	runningResourcePollers atomic.Int64
	resourcePollCount      atomic.Int64
	resourceReplyCount     atomic.Int64
	resourceCleanupCount   atomic.Int64
}

type mdnsWatcher struct {
	servicesCh chan *mdns.ServiceEntry
	stats      mdnsWatcherStats

	sync.Mutex
	services mdnsServiceTracker

	cancelFuncs map[string]context.CancelFunc
	entities    map[string]*mdnsService
}

func newmdnsWatcher() *mdnsWatcher {
	return &mdnsWatcher{
		servicesCh:  make(chan *mdns.ServiceEntry, 1024),
		cancelFuncs: make(map[string]context.CancelFunc),
		entities:    make(map[string]*mdnsService),
	}
}

func (m *mdnsWatcher) handleOneIncomingServiceReply(name string) {
	m.stats.serviceReplyCount.Add(1)

	m.Lock()
	defer m.Unlock()

	// Is this a valid service name?
	if strings.Count(name, ".") != 3 {
		return
	}

	parts := strings.Split(name, ".")

	replaced := m.services.observe(name)
	if !replaced {
		logutil.Verbosef("+ %+s [%s]\n", name, serviceDescriptions[parts[0]])
	}

	if !replaced {
		m.spawnService(parts[0]+"."+parts[1]+".", parts[2])
	}
}

func (m *mdnsWatcher) handleIncomingServiceReplies() {
	for entry := range m.servicesCh {
		m.handleOneIncomingServiceReply(entry.Name)
	}
}

func (m *mdnsWatcher) sendServiceRequests() {
	t := time.NewTicker(20 * time.Second)
	for {
		go m.sendOneServiceRequest()
		<-t.C
	}
}

func (m *mdnsWatcher) sendOneServiceRequest() {
	m.stats.runningServicePollers.Add(1)
	m.stats.servicePollCount.Add(1)
	defer m.stats.runningServicePollers.Add(-1)
	mdns.Query(&mdns.QueryParam{
		Service: "_services._dns-sd._udp",
		Timeout: 60 * time.Second,
		Entries: m.servicesCh,
		Logger:  log.New(io.Discard, "", 0),
	})
}

func (m *mdnsWatcher) removeStaleJunkOnce() {
	m.stats.serviceCleanupCount.Add(1)

	m.Lock()
	defer m.Unlock()

	removed := m.services.removeStaleJunk(45 * time.Minute)
	for _, svc := range removed {
		logutil.Verbosef(" - %s\n", svc)

		parts := strings.Split(svc, ".")
		m.killService(parts[0]+"."+parts[1]+".", parts[2])
	}
}

func (m *mdnsWatcher) removeStaleJunkRepeatedly() {
	t := time.NewTicker(10 * time.Minute)
	for _ = range t.C {
		m.removeStaleJunkOnce()
	}
}

func (m *mdnsWatcher) watchServices() {
	go m.handleIncomingServiceReplies()
	go m.removeStaleJunkRepeatedly()

	m.sendServiceRequests()
}

// Per-service listener gubbins

func (m *mdnsWatcher) spawnService(name string, domain string) {
	// only call this if you already have a lock pls
	logutil.Verbosef("+ spawning %s\n", name)

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFuncs[name] = cancel
	entity := newMDNSService(name, domain, &m.stats)
	m.entities[name] = entity

	go entity.perServiceRunloop(ctx)
}

func (m *mdnsWatcher) killService(name string, domain string) {
	// only call this if you already have a lock pls
	logutil.Verbosef("- killing %s\n", name)

	cancel, ok := m.cancelFuncs[name]
	if !ok {
		return
	}

	cancel()
	delete(m.cancelFuncs, name)
	delete(m.entities, name)
}

type serviceDetails struct {
	Name   string
	Host   string
	AddrV4 string
	AddrV6 string
	Port   int
	Info   string
}

func (s *serviceDetails) resource(serviceName string) Resource {
	return Resource{
		Service: serviceName,
		Name:    s.Name,
		Host:    s.Host,
		AddrV4:  s.AddrV4,
		AddrV6:  s.AddrV6,
		Port:    s.Port,
		Info:    s.Info,
	}
}

func mkServiceDetails(entry *mdns.ServiceEntry) serviceDetails {
	deets := serviceDetails{
		Name: entry.Name,
		Host: entry.Host,
		Port: entry.Port,
		Info: entry.Info,
	}

	if entry.AddrV4 != nil {
		deets.AddrV4 = entry.AddrV4.String()
	}

	if entry.AddrV6 != nil {
		deets.AddrV6 = entry.AddrV6.String()
	}

	return deets
}

type mdnsService struct {
	c     chan *mdns.ServiceEntry
	stats *mdnsWatcherStats

	sync.RWMutex

	name   string
	domain string
	m      map[serviceDetails]time.Time
}

func newMDNSService(name string, domain string, stats *mdnsWatcherStats) *mdnsService {
	return &mdnsService{
		c:      make(chan *mdns.ServiceEntry, 1024),
		stats:  stats,
		name:   name,
		domain: domain,
		m:      make(map[serviceDetails]time.Time),
	}
}

func (m *mdnsService) removeStaleJunkOnce() {
	m.stats.resourceCleanupCount.Add(1)

	m.Lock()
	defer m.Unlock()

	for k, v := range m.m {
		if time.Since(v) > 25*time.Minute {
			logutil.Verbosef("-- %s\n", k.Name)
			delete(m.m, k)
		}
	}
}

func (m *mdnsService) removeStaleJunk(ctx context.Context) {
	t := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.removeStaleJunkOnce()
		}
	}
}

func (m *mdnsService) sendOneRequest(ctx context.Context) {
	m.stats.resourcePollCount.Add(1)
	m.stats.runningResourcePollers.Add(1)
	defer m.stats.runningResourcePollers.Add(-1)
	mdns.QueryContext(ctx, &mdns.QueryParam{
		Service: m.name,
		Domain:  m.domain,
		Timeout: 60 * time.Second,
		Entries: m.c,
		Logger:  log.New(io.Discard, "", 0),
	})
}

func (m *mdnsService) sendRequests(ctx context.Context) {
	m.sendOneRequest(ctx)

	t := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			go m.sendOneRequest(ctx)
		}
	}
}

func (m *mdnsService) handleOneReply(e *mdns.ServiceEntry) {
	m.stats.resourceReplyCount.Add(1)

	m.Lock()
	defer m.Unlock()

	// filter on suffix
	suffix := m.name + m.domain + "."
	if !strings.HasSuffix(e.Name, suffix) {
		return
	}

	ds := mkServiceDetails(e)

	_, exists := m.m[ds]
	if !exists {
		logutil.Verbosef("++ (%s) %s @ %v\n", m.name, e.Name, e.TTL)
	}

	m.m[ds] = time.Now()
}

func (m *mdnsService) handleReplies(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-m.c:
			m.handleOneReply(e)
		}
	}
}

func (m *mdnsService) perServiceRunloop(ctx context.Context) {
	go m.removeStaleJunk(ctx)
	go m.handleReplies(ctx)
	m.sendRequests(ctx)
}

// Metrics generation

func (m *mdnsServiceTracker) fillMetrics(metrics *Metrics) {
	metrics.ServiceLastSeen = make(map[Service]int64)

	m.RLock()
	defer m.RUnlock()

	for k, v := range m.m {
		parts := strings.Split(k, ".")
		if len(parts) < 3 {
			continue
		}

		svc := Service{
			Name: k,
			Desc: serviceDescriptions[parts[0]],
		}
		metrics.ServiceLastSeen[svc] = v.Unix()
	}
}

func (m *mdnsService) fillMetrics(metrics *Metrics) {
	m.RLock()
	defer m.RUnlock()

	if metrics.ResourceLastSeen == nil {
		metrics.ResourceLastSeen = make(map[Resource]int64)
	}

	for k, v := range m.m {
		res := k.resource(m.name + m.domain + ".")
		metrics.ResourceLastSeen[res] = v.Unix()
	}
}

func (m *mdnsWatcherStats) fillMetrics(metrics *Metrics) {
	metrics.RunningServicePollers = m.runningServicePollers.Load()
	metrics.ServicePollCount = m.servicePollCount.Load()
	metrics.ServiceReplyCount = m.serviceReplyCount.Load()
	metrics.ServiceCleanupCount = m.serviceCleanupCount.Load()

	metrics.RunningResourcePollers = m.runningResourcePollers.Load()
	metrics.ResourcePollCount = m.resourcePollCount.Load()
	metrics.ResourceReplyCount = m.resourceReplyCount.Load()
	metrics.ResourceCleanupCount = m.resourceCleanupCount.Load()
}

func (m *mdnsWatcher) metrics() Metrics {
	m.Lock()
	defer m.Unlock()

	metrics := Metrics{}
	m.services.fillMetrics(&metrics)

	for _, e := range m.entities {
		e.fillMetrics(&metrics)
	}

	m.stats.fillMetrics(&metrics)

	return metrics
}
