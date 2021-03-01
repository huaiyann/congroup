package congroup

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestConGroup_multi(t *testing.T) {
	wg := New()

	var execCnt int64
	fErr1 := func() error {
		atomic.AddInt64(&execCnt, 1)
		return errors.New("f_err_1")
	}
	fErr2 := func() error {
		atomic.AddInt64(&execCnt, 1)
		return errors.New("f_err_2")
	}
	fNoErr := func() error {
		atomic.AddInt64(&execCnt, 1)
		return nil
	}
	fPanic := func() error {
		atomic.AddInt64(&execCnt, 1)
		panic("f_panic")
	}

	tests := []struct {
		name    string
		funcs   []func() error
		wantErr bool
	}{
		{
			name:    "no_err",
			funcs:   []func() error{fNoErr},
			wantErr: false,
		},
		{
			name:    "no_err_multi",
			funcs:   []func() error{fNoErr, fNoErr},
			wantErr: false,
		},
		{
			name:    "has_err",
			funcs:   []func() error{fNoErr, fErr1, fErr2},
			wantErr: true,
		},
		{
			name:    "has_panic",
			funcs:   []func() error{fNoErr, fPanic},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, f := range tt.funcs {
				f := f
				wg.Add(f)
			}
			if err := wg.Wait(); (err != nil) != tt.wantErr {
				t.Errorf("ConGroup.Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
			if executed := atomic.SwapInt64(&execCnt, 0); executed != int64(len(tt.funcs)) {
				t.Errorf("want %d func executed, but %d", len(tt.funcs), executed)
			}
		})
	}
}
