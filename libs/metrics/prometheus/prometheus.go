/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/07/05        Yang Chenglong
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type LabelValues []string

func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}

type Counter struct {
	cv  *prometheus.CounterVec
	lvs LabelValues
}

func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return NewCounter(cv)
}

func NewCounter(cv *prometheus.CounterVec) *Counter {
	return &Counter{
		cv: cv,
	}
}

func (c *Counter) With(labelValues ...string) *Counter {
	return &Counter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}

func (c *Counter) Add(delta float64) {
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}

type Gauge struct {
	gv  *prometheus.GaugeVec
	lvs LabelValues
}

func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return NewGauge(gv)
}

func NewGauge(gv *prometheus.GaugeVec) *Gauge {
	return &Gauge{
		gv: gv,
	}
}

func (g *Gauge) With(labelValues ...string) *Gauge {
	return &Gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}

func (g *Gauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}

func (g *Gauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}

type Summary struct {
	sv  *prometheus.SummaryVec
	lvs LabelValues
}

func NewSummaryFrom(opts prometheus.SummaryOpts, labelNames []string) *Summary {
	sv := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(sv)
	return NewSummary(sv)
}

func NewSummary(sv *prometheus.SummaryVec) *Summary {
	return &Summary{
		sv: sv,
	}
}

func (s *Summary) With(labelValues ...string) *Summary {
	return &Summary{
		sv:  s.sv,
		lvs: s.lvs.With(labelValues...),
	}
}

func (s *Summary) Observe(value float64) {
	s.sv.With(makeLabels(s.lvs...)).Observe(value)
}

type Histogram struct {
	hv  *prometheus.HistogramVec
	lvs LabelValues
}

func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return NewHistogram(hv)
}

func NewHistogram(hv *prometheus.HistogramVec) *Histogram {
	return &Histogram{
		hv: hv,
	}
}

func (h *Histogram) With(labelValues ...string) *Histogram {
	return &Histogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}

func (h *Histogram) Observe(value float64) {
	h.hv.With(makeLabels(h.lvs...)).Observe(value)
}

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
