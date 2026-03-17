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

#define for0(i, n) for (int i = 0; i < (int)(n); ++i)
#define for1(i, n) for (int i = 1; i <= (int)(n); ++i)
#define forc(i, l, r) for (int i = (int)(l); i <= (int)(r); ++i)
#define forr0(i, n) for (int i = (int)(n) - 1; i >= 0; --i)
#define forr1(i, n) for (int i = (int)(n); i >= 1; --i)
#define each(x, a) for (auto &x : a)

#define pb push_back
#define fi first
#define se second
#define eb emplace_back
#define ef emplace_front
#define em emplace
#define fr front()
#define bk back()

#define bpc __builtin_popcount
#define bpcll __builtin_popcountll
#define clz __builtin_clz
#define clzll __builtin_clzll
#define ctzll __builtin_ctzll
#define ctz __builtin_ctz
#define sqrt __builtin_sqrt
#define abs __builtin_abs
#define memset __builtin_memset
#define memcpy __builtin_memcpy

#define all(x) (x).begin(), (x).end()
#define rall(x) (x).rbegin(), (x).rend()

#define present(c, x) ((c).find(x) != (c).end())
#define cpresent(c, x) (find(all(c), x) != c.end())

#define wne(c) while (!((c).empty()))

#define sz(a) int((a).size())

using namespace std;

using ll = long long;
using db = double;
using ld = long double;
using ul = unsigned long;
using ull = unsigned long long;

using vi = vector<int>;
using vvi = vector<vi>;
using pi = pair<int, int>;
using pll = pair<ll, ll>;
using pdb = pair<db, db>;
using vpi = vector<pi>;
using vc = vector<char>;
using vdb = vector<db>;
using vs = vector<string>;
using vll = vector<ll>;
using vvll = vector<vll>;
using vb = vector<bool>;
using si = unordered_set<int>;
using mi = unordered_map<int, int>;
using sc = unordered_set<char>;
using pqi = priority_queue<int>;
using pqpi = priority_queue<pi>;
using pqll = priority_queue<ll>;
using pqpll = priority_queue<pll>;

template <typename T> void print(vector<T> x) {
  for (auto i : x)
    cout << i << ' ';
  cout << "\n";
}
template <typename T> void print(set<T> x) {
  for (auto i : x)
    cout << i << ' ';
  cout << "\n";
}
template <typename T> void print(unordered_set<T> x) {
  for (auto i : x)
    cout << i << ' ';
  cout << "\n";
}
template <typename T> void print(T &&x) { cout << x << "\n"; }
template <typename T, typename... S> void print(T &&x, S &&...y) {
  cout << x << ' ';
  print(y...);
}

template <typename T> istream &operator>>(istream &i, vector<T> &vec) {
  for (auto &x : vec)
    i >> x;
  return i;
}

vvi read_graph(int n, int m, int base = 1) {
  vvi adj(n);
  for (int i = 0, u, v; i < m; ++i) {
    cin >> u >> v, u -= base, v -= base;
    adj[u].pb(v), adj[v].pb(u);
  }
  return adj;
}

vvi read_tree(int n, int base = 1) { return read_graph(n, n - 1, base); }

template <typename T, typename S>
pair<T, S> operator+(const pair<T, S> &a, const pair<T, S> &b) {
  return {a.first + b.first, a.second + b.second};
}

const int MOD = 1000000007;

vector<ll> fact, invfact;

ll modpow(ll a, ll e) {
  ll r = 1;
  while (e) {
    if (e & 1)
      r = r * a % MOD;
    a = a * a % MOD;
    e >>= 1;
  }
  return r;
}

void init_ncr(int N) {
  fact.resize(N + 1);
  invfact.resize(N + 1);

  fact[0] = 1;
  for (int i = 1; i <= N; i++)
    fact[i] = fact[i - 1] * i % MOD;

  invfact[N] = modpow(fact[N], MOD - 2);
  for (int i = N; i > 0; i--)
    invfact[i - 1] = invfact[i] * i % MOD;
}

ll ncr(int n, int r) {
  if (r < 0 || r > n)
    return 0;
  return fact[n] * invfact[r] % MOD * invfact[n - r] % MOD;
}

int fastlog(int x) { return 64 - __builtin_clzll(x) - 1; }

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
