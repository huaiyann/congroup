package congroup

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestConGroupErrs(t *testing.T) {
	ctx := context.Background()
	cg := New(ctx)
	cg.Add(func(ctx context.Context) error { return nil })
	cg.Add(func(ctx context.Context) error { return errors.New("fake 2 error") })
	cg.Add(func(ctx context.Context) error { return errors.New("fake 1 error") })
	err := cg.Wait()
	if err == nil {
		t.Fatal("want error but nil")
	}
	lines := strings.Split(err.Error(), "\n")
	if wantLine, gotLine := 7, len(lines); wantLine != gotLine {
		t.Fatalf("want error has %d line, but %d", wantLine, gotLine)
	}
	if wantPrefix, gotPrefix := "2 errors occurred", lines[0]; wantPrefix != gotPrefix {
		t.Fatalf("want prefix %s, but %s", wantPrefix, gotPrefix)
	}
}

func TestConGroupPanics(t *testing.T) {
	ctx := context.Background()
	cg := New(ctx)
	cg.Add(func(ctx context.Context) error { return nil })
	cg.Add(func(ctx context.Context) error { panic("panic reason 1") })
	cg.Add(func(ctx context.Context) error { panic("panic reason 2") })

	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("want panic but not")
		}
		info := err.(string)

		wantReg := `NO\.[0-9]: panic reason 1`
		if reg := regexp.MustCompile(wantReg); !reg.MatchString(info) {
			t.Fatalf("want match regex %s, but not", wantReg)
		}
		wantReg = `NO\.[0-9]: panic reason 2`
		if reg := regexp.MustCompile(wantReg); !reg.MatchString(info) {
			t.Fatalf("want match regex %s, but not", wantReg)
		}
		lines := strings.Split(info, "\n")
		if wantPrefix, gotPrefix := "2 panics occurred", lines[0]; wantPrefix != gotPrefix {
			log.Fatalf("want prefix %s, but %s", wantPrefix, gotPrefix)
		}
	}()

	cg.Wait()
	t.Fatal("should when wait")
}

func TestConGroupResult(t *testing.T) {
	var got, want int64
	want = 100000
	ctx := context.Background()
	cg := New(ctx)
	for i := int64(0); i < want; i++ {
		cg.Add(func(ctx context.Context) error {
			atomic.AddInt64(&got, 1)
			return nil
		})
	}
	err := cg.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("want result %d, but %d", want, got)
	}
}

func BenchmarkConGroupAdd(b *testing.B) {
	ctx := context.Background()
	cg := New(ctx)
	for i := 0; i < b.N; i++ {
		cg.Add(func(ctx context.Context) error { return nil })
	}
	err := cg.Wait()
	if err != nil {
		b.Fatal(err)
	}
}

func TestConGroupContextTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	cg := New(ctx)
	cg.Add(func(ctx context.Context) error { return nil })
	cg.Add(func(ctx context.Context) error { return nil })
	cg.Add(func(ctx context.Context) error {
		<-time.After(time.Millisecond * 2)
		return nil
	})
	err := cg.Wait()
	if err == nil {
		t.Fatal("want error but nil")
	}
	lines := strings.Split(err.Error(), "\n")
	if wantLine, gotLine := 4, len(lines); wantLine != gotLine {
		t.Fatalf("want error has %d line, but %d", wantLine, gotLine)
	}
	if wantErr, gotErr := "NO.1: context deadline exceeded", lines[2]; wantErr != gotErr {
		t.Fatalf("want error %s, but %s", wantErr, gotErr)
	}
}
