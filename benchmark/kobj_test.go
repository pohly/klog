/*
Copyright 2021 The Kubernetes Authors.

SPDX-License-Identifier: Apache-2.0
*/

package example

import (
	"flag"
	"io"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type kmeta struct {
	name, namespace string
}

func (k kmeta) GetName() string {
	return k.name
}

func (k kmeta) GetNamespace() string {
	return k.namespace
}

var _ klog.KMetadata = kmeta{}

// var obj = kmeta{name: "some-fake-name", namespace: "kube-system"}
var obj = v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "some-fake-name", Namespace: "kube-system"}}

var result string

func BenchmarkString(b *testing.B) {
	var s string
	for n := 0; n < b.N; n++ {
		s = klog.KObj(&obj).String()
	}
	result = s
}

func BenchmarkKlog(b *testing.B) {
	klog.SetOutput(io.Discard)
	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	flags.Set("logtostderr", "false")
	for n := 0; n < b.N; n++ {
		klog.InfoS("printed", "obj", klog.KObj(&obj))
	}
}

func BenchmarkKlogSkipped(b *testing.B) {
	for n := 0; n < b.N; n++ {
		klog.V(10).InfoS("skipped", "obj", klog.KObj(&obj))
	}
}

func BenchmarkDiscard(b *testing.B) {
	log := logr.Discard()
	for n := 0; n < b.N; n++ {
		log.Info("skipped", klog.KObj(&obj))
	}
}

func BenchmarkZap(b *testing.B) {
	w := zapcore.AddSync(io.Discard)
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		CallerKey:  "caller",
		TimeKey:    "ts",
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(w), zap.DebugLevel)
	l := zap.New(core, zap.WithCaller(false))
	log := zapr.NewLogger(l)

	for n := 0; n < b.N; n++ {
		log.Info("printed", "obj", klog.KObj(&obj))
	}
}

func BenchmarkZapSkipped(b *testing.B) {
	w := zapcore.AddSync(io.Discard)
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		CallerKey:  "caller",
		TimeKey:    "ts",
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(w), zap.DebugLevel)
	l := zap.New(core, zap.WithCaller(false))
	log := zapr.NewLogger(l)

	for n := 0; n < b.N; n++ {
		log.V(10).Info("skipped", "obj", klog.KObj(&obj))
	}
}
