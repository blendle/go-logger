package stackdriver

import "go.uber.org/zap/zapcore"

type label struct {
	Key   string
	Value string
}

type labels []*label

func (l labels) Clone() labels {
	return l
}

func (l labels) MarshalLogObject(e zapcore.ObjectEncoder) error {
	for i := range l {
		e.AddString(l[i].Key, l[i].Value)
	}
	return nil
}
