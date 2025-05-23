// Copyright 2024 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package selfmonitor

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/alibaba/ilogtail/pkg/helper/pool"
	"github.com/alibaba/ilogtail/pkg/protocol"
)

const (
	defaultTagValue = "-"
)

var (
	DefaultCacheFactory = NewMapCache
)

// SetMetricVectorCacheFactory allows users to set the cache factory for the metric vector, like Prometheus SDK.
func SetMetricVectorCacheFactory(factory func(MetricSet) MetricVectorCache) {
	DefaultCacheFactory = factory
}

type (
	CumulativeCounterMetricVector = MetricVector[CounterMetric]
	AverageMetricVector           = MetricVector[CounterMetric]
	MaxMetricVector               = MetricVector[GaugeMetric]
	CounterMetricVector           = MetricVector[CounterMetric]
	GaugeMetricVector             = MetricVector[GaugeMetric]
	LatencyMetricVector           = MetricVector[LatencyMetric]
	StringMetricVector            = MetricVector[StringMetric]
)

// Deprecated: use NewCounterMetricVector instead.
// NewCumulativeCounterMetricVector creates a new CounterMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewCumulativeCounterMetricVector(metricName string, constLabels map[string]string, labelNames []string) CumulativeCounterMetricVector {
	return NewMetricVector[CounterMetric](metricName, CumulativeCounterType, constLabels, labelNames)
}

// NewCounterMetricVector creates a new DeltaMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewCounterMetricVector(metricName string, constLabels map[string]string, labelNames []string) CounterMetricVector {
	return NewMetricVector[CounterMetric](metricName, CounterType, constLabels, labelNames)
}

// NewAverageMetricVector creates a new AverageMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewAverageMetricVector(metricName string, constLabels map[string]string, labelNames []string) AverageMetricVector {
	return NewMetricVector[CounterMetric](metricName, AverageType, constLabels, labelNames)
}

// NewMaxMetricVector creates a new MaxMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewMaxMetricVector(metricName string, constLabels map[string]string, labelNames []string) MaxMetricVector {
	return NewMetricVector[GaugeMetric](metricName, MaxType, constLabels, labelNames)
}

// NewGaugeMetricVector creates a new GaugeMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewGaugeMetricVector(metricName string, constLabels map[string]string, labelNames []string) GaugeMetricVector {
	return NewMetricVector[GaugeMetric](metricName, GaugeType, constLabels, labelNames)
}

// NewStringMetricVector creates a new StringMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewStringMetricVector(metricName string, constLabels map[string]string, labelNames []string) StringMetricVector {
	return NewMetricVector[StringMetric](metricName, StringType, constLabels, labelNames)
}

// NewLatencyMetricVector creates a new LatencyMetricVector.
// Note that MetricVector doesn't expose Collect API by default. Plugins Developers should be careful to collect metrics manually.
func NewLatencyMetricVector(metricName string, constLabels map[string]string, labelNames []string) LatencyMetricVector {
	return NewMetricVector[LatencyMetric](metricName, LatencyType, constLabels, labelNames)
}

// NewCumulativeCounterMetricVectorAndRegister creates a new CounterMetricVector and register it to the MetricsRecord.
func NewCumulativeCounterMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) CumulativeCounterMetricVector {
	v := NewMetricVector[CounterMetric](metricName, CumulativeCounterType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewAverageMetricVectorAndRegister creates a new AverageMetricVector and register it to the MetricsRecord.
func NewAverageMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) AverageMetricVector {
	v := NewMetricVector[CounterMetric](metricName, AverageType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewCounterMetricVectorAndRegister creates a new DeltaMetricVector and register it to the MetricsRecord.
func NewCounterMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) CounterMetricVector {
	v := NewMetricVector[CounterMetric](metricName, CounterType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewGaugeMetricVectorAndRegister creates a new GaugeMetricVector and register it to the MetricsRecord.
func NewGaugeMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) GaugeMetricVector {
	v := NewMetricVector[GaugeMetric](metricName, GaugeType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewLatencyMetricVectorAndRegister creates a new LatencyMetricVector and register it to the MetricsRecord.
func NewLatencyMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) LatencyMetricVector {
	v := NewMetricVector[LatencyMetric](metricName, LatencyType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewStringMetricVectorAndRegister creates a new StringMetricVector and register it to the MetricsRecord.
func NewStringMetricVectorAndRegister(mr *MetricsRecord, metricName string, constLabels map[string]string, labelNames []string) StringMetricVector {
	v := NewMetricVector[StringMetric](metricName, StringType, constLabels, labelNames)
	mr.RegisterMetricCollector(v)
	return v
}

// NewCumulativeCounterMetric creates a new CounterMetric.
func NewCumulativeCounterMetric(n string, lables ...*protocol.Log_Content) CounterMetric {
	return NewCumulativeCounterMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewAverageMetric creates a new AverageMetric.
func NewAverageMetric(n string, lables ...*protocol.Log_Content) CounterMetric {
	return NewAverageMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewCounterMetric creates a new DeltaMetric.
func NewCounterMetric(n string, lables ...*protocol.Log_Content) CounterMetric {
	return NewCounterMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewGaugeMetric creates a new GaugeMetric.
func NewGaugeMetric(n string, lables ...*protocol.Log_Content) GaugeMetric {
	return NewGaugeMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewStringMetric creates a new StringMetric.
func NewStringMetric(n string, lables ...*protocol.Log_Content) StringMetric {
	return NewStringMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewLatencyMetric creates a new LatencyMetric.
func NewLatencyMetric(n string, lables ...*protocol.Log_Content) LatencyMetric {
	return NewLatencyMetricVector(n, convertLabels(lables), nil).WithLabels()
}

// NewCumulativeCounterMetricAndRegister creates a new CounterMetric and register it's metricVector to the MetricsRecord.
func NewCumulativeCounterMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) CounterMetric {
	mv := NewCumulativeCounterMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewCounterMetricAndRegister creates a new DeltaMetric and register it's metricVector to the MetricsRecord.
func NewCounterMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) CounterMetric {
	mv := NewCounterMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewAverageMetricAndRegister creates a new AverageMetric and register it's metricVector to the MetricsRecord.
func NewAverageMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) CounterMetric {
	mv := NewAverageMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewMaxMetricAndRegister creates a new MaxMetric and register it's metricVector to the MetricsRecord.
func NewMaxMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) GaugeMetric {
	mv := NewMaxMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewGaugeMetricAndRegister creates a new GaugeMetric and register it's metricVector to the MetricsRecord.
func NewGaugeMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) GaugeMetric {
	mv := NewGaugeMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewLatencyMetricAndRegister creates a new LatencyMetric and register it's metricVector to the MetricsRecord.
func NewLatencyMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) LatencyMetric {
	mv := NewLatencyMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

// NewStringMetricAndRegister creates a new StringMetric and register it's metricVector to the MetricsRecord.
func NewStringMetricAndRegister(c *MetricsRecord, n string, lables ...*protocol.Log_Content) StringMetric {
	mv := NewStringMetricVector(n, convertLabels(lables), nil)
	c.RegisterMetricCollector(mv.(MetricCollector))
	return mv.WithLabels()
}

var (
	_ MetricCollector             = (*MetricVectorImpl[CounterMetric])(nil)
	_ MetricSet                   = (*MetricVectorImpl[StringMetric])(nil)
	_ MetricVector[CounterMetric] = (*MetricVectorImpl[CounterMetric])(nil)
	_ MetricVector[GaugeMetric]   = (*MetricVectorImpl[GaugeMetric])(nil)
	_ MetricVector[LatencyMetric] = (*MetricVectorImpl[LatencyMetric])(nil)
	_ MetricVector[StringMetric]  = (*MetricVectorImpl[StringMetric])(nil)
)

type MetricVectorAndCollector[T Metric] interface {
	MetricVector[T]
	MetricCollector
}

type MetricVectorImpl[T Metric] struct {
	*metricVector
}

// NewMetricVector creates a new MetricVector.
// It returns a MetricVectorAndCollector, which is a MetricVector and a MetricCollector.
// For plugin developers, they should use MetricVector APIs to create metrics.
// For agent itself, it uses MetricCollector APIs to collect metrics.
func NewMetricVector[T Metric](metricName string, metricType SelfMetricType, constLabels map[string]string, labelNames []string) MetricVectorAndCollector[T] {
	return &MetricVectorImpl[T]{
		metricVector: newMetricVector(metricName, metricType, constLabels, labelNames),
	}
}

func (m *MetricVectorImpl[T]) WithLabels(labels ...LabelPair) T {
	return m.metricVector.WithLabels(labels...).(T)
}

// MetricVectorCache is a cache for MetricVector.
type MetricVectorCache interface {
	// return a metric with the given label values.
	// Note that the label values are sorted according to the label keys in MetricSet.
	WithLabelValues([]string) Metric

	MetricCollector
}

type metricVector struct {
	name        string // metric name
	metricType  SelfMetricType
	constLabels []LabelPair // constLabels is the labels that are not changed when the metric is created.
	labelKeys   []string    // labelNames is the names of the labels. The values of the labels can be changed.

	indexPool pool.GenericPool[string] // index is []string, which is sorted according to labelNames.
	cache     MetricVectorCache        // collector is a map[string]Metric, key is the index of the metric.
}

func newMetricVector(
	metricName string,
	metricType SelfMetricType,
	constLabels map[string]string,
	labelNames []string,
) *metricVector {
	mv := &metricVector{
		name:       metricName,
		metricType: metricType,
		labelKeys:  labelNames,
		indexPool:  pool.NewGenericPool(func() []string { return make([]string, 0, 10) }),
	}

	for k, v := range constLabels {
		mv.constLabels = append(mv.constLabels, LabelPair{Key: k, Value: v})
	}

	mv.cache = DefaultCacheFactory(mv)
	return mv
}

func (v *metricVector) Name() string {
	return v.name
}

func (v *metricVector) Type() SelfMetricType {
	return v.metricType
}

func (v *metricVector) ConstLabels() []LabelPair {
	return v.constLabels
}

func (v *metricVector) LabelKeys() []string {
	return v.labelKeys
}

func (v *metricVector) WithLabels(labels ...LabelPair) Metric {
	labelValues, err := v.buildLabelValues(labels)
	if err != nil {
		return newErrorMetric(v.metricType, err)
	}
	defer v.indexPool.Put(labelValues)
	return v.cache.WithLabelValues(*labelValues)
}

func (v *metricVector) Collect() []Metric {
	return v.cache.Collect()
}

// buildLabelValues return the index
func (v *metricVector) buildLabelValues(labels []LabelPair) (*[]string, error) {
	if len(labels) > len(v.labelKeys) {
		return nil, fmt.Errorf("too many labels, expected %d, got %d. defined labels: %v",
			len(v.labelKeys), len(labels), v.labelKeys)
	}

	index := v.indexPool.Get()
	for range v.labelKeys {
		*index = append(*index, defaultTagValue)
	}

	for d, tag := range labels {
		if v.labelKeys[d] == tag.Key { // fast path
			(*index)[d] = tag.Value
		} else {
			err := v.slowConstructIndex(index, tag)
			if err != nil {
				v.indexPool.Put(index)
				return nil, err
			}
		}
	}

	return index, nil
}

func (v *metricVector) slowConstructIndex(index *[]string, tag LabelPair) error {
	for i, tagName := range v.labelKeys {
		if tagName == tag.Key {
			(*index)[i] = tag.Value
			return nil
		}
	}
	return fmt.Errorf("undefined label: %s in %v", tag.Key, v.labelKeys)
}

type MapCache struct {
	MetricSet
	bytesPool pool.GenericPool[byte]
	sync.Map
}

func NewMapCache(metricSet MetricSet) MetricVectorCache {
	return &MapCache{
		MetricSet: metricSet,
		bytesPool: pool.NewGenericPool(func() []byte { return make([]byte, 0, 128) }),
	}
}

func (v *MapCache) WithLabelValues(labelValues []string) Metric {
	buffer := v.bytesPool.Get()
	for _, tagValue := range labelValues {
		*buffer = append(*buffer, '|')
		*buffer = append(*buffer, tagValue...)
	}

	/* #nosec G103 */
	k := *(*string)(unsafe.Pointer(buffer))
	acV, loaded := v.Load(k)
	if loaded {
		metric := acV.(Metric)
		v.bytesPool.Put(buffer)
		return metric
	}

	newMetric := newMetric(v.Type(), v, labelValues)
	acV, loaded = v.LoadOrStore(k, newMetric)
	if loaded {
		v.bytesPool.Put(buffer)
	}
	return acV.(Metric)
}

func (v *MapCache) Collect() []Metric {
	res := make([]Metric, 0, 10)
	v.Range(func(key, value interface{}) bool {
		res = append(res, value.(Metric))
		return true
	})
	return res
}

func convertLabels(labels []*protocol.Log_Content) map[string]string {
	if len(labels) == 0 {
		return nil
	}

	l := make(map[string]string)
	for _, label := range labels {
		l[label.Key] = label.Value
	}

	return l
}
