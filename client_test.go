package etumag1matrix

import "testing"

var (
	bad_addrs = []string{
		"",
		"localhost",
		"127.0.0.1",
	}
	good_addrs = []string{
		"localhost:80",
		"127.0.0.1:80",
		":9090",
	}
	much_good_addrs = []string{
		"localhost:11",
		"127.0.0.1:12",
		":13",
		"127.0.0.1:14",
		"127.0.0.1:15",
		"192.168.1.1:16",
		"192.168.1.1:17",
		"192.168.1.1:18",
		"192.168.1.1:19",
		"192.168.1.1:20",
	}
)

func TestNewClient(t *testing.T) {
	var err error
	_, err = NewClient(bad_addrs...)
	if err == nil {
		t.Error("NewClient fail:", err)
	}

	_, err = NewClient(good_addrs...)
	if err != nil {
		t.Error("NewClient fail:", err)
	}

	_, err = NewClient(much_good_addrs...)
	if err != nil {
		t.Error("NewClient fail:", err)
	}
}

func TestGetAddrs(t *testing.T) {
	client, err := NewClient(much_good_addrs...)
	if err != nil {
		t.Error("NewClient fail:", err)
	}
	if len(client.getAddrs(0)) > 0 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(100)) != len(client.addrs) {
		t.Error("getAddrs fail")
	}

	if client.cursor != 0 {
		t.Error("cursor error")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 3 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 6 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 9 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(1)) != 1 || client.cursor != 0 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(4)) != 4 || client.cursor != 4 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 7 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(2)) != 2 || client.cursor != 9 {
		t.Error("getAddrs fail")
	}

	if len(client.getAddrs(3)) != 3 || client.cursor != 2 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 5 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 8 {
		t.Error("getAddrs fail")
	}
	if len(client.getAddrs(3)) != 3 || client.cursor != 1 {
		t.Error("getAddrs fail")
	}
}

var matrixFromSlice_ = NewMatrixFromSlice(3, 3, []float64{
	1, 2, 3, 4, // 5, 6,
	7, 8, 9, 10, // 11, 12,
	13, 14, 15, 16, // 17, 18,
	19, 20, 21, 22, //23, 24,
	// 25, 26, 27, 28, 29, 30,
})

func TestMain(t *testing.T) {
	client, _ := NewClient("127.0.0.1:8080")
	t.Error(client.MultiplyMatrix(matrixFromSlice_, matrixFromSlice_))
}
