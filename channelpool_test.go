package proxy

import "testing"

func Test_chanPool(t *testing.T) {
	factoty := func(addr string) (*ReverseProxy, error) {
		p := NewReverseProxy(addr)
		return p, nil
	}

	pool, err := NewChanPool(5, 100, factoty)
	if err != nil {
		t.Fatalf("could not make chan pool: %v", err)
	}

	t.Logf("len of pool is %d", pool.Len())

	p, err := pool.Get("localhost:8080")
	if err != nil {
		t.Fatalf("could not make chan pool: %v", err)
	}

	if p == nil {
		t.Fatalf("could not get one proxy form pool, proxy is nil")
	}

	t.Logf("proxy addr: %v and addr is: %s", p, p.getClient().Addr)
}

func BenchmarkNewReverseProxyWithPool(b *testing.B) {
	b.StopTimer()
	pool, err := NewChanPool(10, 100, func(addr string) (*ReverseProxy, error) {
		return NewReverseProxy(addr), nil
	})
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		proxy, err := pool.Get("locahost:8080")
		if err != nil {
			b.Fatal(err)
		}
		if proxy == nil {
			b.Fatalf("could not get from pool, proxy is nil")
		}
		if proxy.getClient() == nil {
			b.Fatalf("could not get from pool, client is nil")
		}
	}
}
