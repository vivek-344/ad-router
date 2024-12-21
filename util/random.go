package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int32) int32 {
	return min + r.Int31n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomBool() bool {
	return r.Intn(2) == 0
}

func RandomCid() string {
	return RandomString(int(RandomInt(6, 8)))
}

func RandomName() string {
	return RandomString(int(RandomInt(7, 10)))
}

func RandomImg() string {
	return "https://example.com"
}

func RandomCta() string {
	cta := []string{"Download", "Install", "Get", "Play"}
	n := len(cta)
	return cta[r.Intn(n)]
}

func RandomAppID() string {
	var sb strings.Builder
	n := int(RandomInt(1, 10))

	for i := 1; i <= n; i++ {
		str := "com." + RandomString(int(RandomInt(5, 10))) + "." + RandomString(int(RandomInt(5, 10)))
		sb.WriteString(str)
		if i < n {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func RandomOs() string {
	var sb strings.Builder
	os := []string{"Android", "iOS", "Web"}

	for i := range os {
		j := rand.Intn(i + 1)
		os[i], os[j] = os[j], os[i]
	}

	n := int(RandomInt(1, 3))

	for i := 0; i < n; i++ {
		sb.WriteString(os[i])
		if i < n-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func RandomCountry() string {
	var sb strings.Builder
	countries := []string{"Russia", "Canada", "China", "United States", "Brazil", "Australia", "India", "Argentina", "Kazakhstan", "Algeria"}

	for i := range countries {
		j := rand.Intn(i + 1)
		countries[i], countries[j] = countries[j], countries[i]
	}

	n := int(RandomInt(1, 10))

	for i := 0; i < n; i++ {
		sb.WriteString(countries[i])
		if i < n-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func RandomRule() string {
	rule := []string{"include", "exclude"}
	n := len(rule)
	return rule[r.Intn(n)]
}
