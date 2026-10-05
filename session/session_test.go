package session

import "testing"

func TestOpenReturnsToken(t *testing.T) {
	m := &Manager{GraceMs: 1000}
	got := m.Open("p1", 0)
	if got.Token == "" || got.State != "active" || got.Version != 1 {
		t.Fatalf("开局会话不对：%+v", got)
	}
}

func TestResumeWithCurrentVersion(t *testing.T) {
	m := &Manager{GraceMs: 1000}
	opened := m.Open("p1", 0)
	snap, reason := m.Resume("r1", opened.Token, opened.Version, 500)
	if reason != "" || snap.Player != "p1" {
		t.Fatalf("正常续局应成功：%+v %q", snap, reason)
	}
}

func TestResumeUnknownToken(t *testing.T) {
	m := &Manager{GraceMs: 1000}
	if _, reason := m.Resume("r1", "nope", 1, 10); reason != "unknown" {
		t.Fatalf("未知令牌应拒绝，实际 %q", reason)
	}
}

func TestTwoPlayersGetDifferentTokens(t *testing.T) {
	m := &Manager{GraceMs: 1000}
	a := m.Open("p1", 0)
	b := m.Open("p2", 0)
	if a.Token == b.Token {
		t.Fatalf("两个玩家的令牌不该相同：%q", a.Token)
	}
}
