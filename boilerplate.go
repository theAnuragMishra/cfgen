package main

import (
	"fmt"
	"time"
)

func generateBoilerplate(name string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(`/**
 *    author: %s
 *    created: %s
 **/
#include <bits/stdc++.h>

using namespace std;

void solve() {}

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int t = 1;
  cin >> t;
  while (t--) {
    solve();
  }
}`, name, now)
}
