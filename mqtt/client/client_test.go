package client

import "testing"

func Test_match(t *testing.T) {
	testCases := []struct {
		in struct {
			topicName   string
			topicFilter string
		}
		out bool
	}{
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/score/wimbledon",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/player1/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis2/",
				topicFilter: "+/tennis/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport//player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/test/player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/player1",
				topicFilter: "sport/+/player1",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "+/monitor/Clients",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients/test",
				topicFilter: "$SYS/monitor/+",
			},
			out: false,
		},
	}

	for _, tc := range testCases {
		if got := match(tc.in.topicName, tc.in.topicFilter); got != tc.out {
			t.Errorf("got => %t, want => %t", got, tc.out)
		}
	}
}
