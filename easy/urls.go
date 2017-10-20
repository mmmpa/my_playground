package easy

import "fmt"

func GenUrls(n int) []string {
	urls := make([]string, 0)

	for i := 1; i <= n; i++ {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}
	return urls
}