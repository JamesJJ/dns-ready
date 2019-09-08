package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jamiealquiza/envy"
	"net"
	"time"
)

type config struct {
	host        *string
	timeout     *int
	retries     *int
	pause       *int64
	verbose     *bool
	silent      *bool
	acceptempty *bool
}

func main() {

	gracefulStop()

	runTime := time.Now().UnixNano()

	conf := config{
		flag.String("host", "kube-dns.kube-system.svc.cluster.local", "hostname to resolve"),
		flag.Int("timeout", 1200, "Timeout in milliseconds for each lookup"),
		flag.Int("retries", 30, "Maximum number of attempts to resolve hostname"),
		flag.Int64("pause", 800, "Number of milliseconds to pause between attempts"),
		flag.Bool("verbose", false, "Show process steps"),
		flag.Bool("silent", false, "Do not print anything to the console"),
		flag.Bool("acceptempty", false, "Accept a DNS response with no IP addresses"),
	}
	envy.Parse("DNSREADY")
	flag.Parse()

	totalAttempts := *conf.retries
	for i := 1; i <= *conf.retries; i++ {
		ips, err := lookupWithTimeout(conf.host, conf.timeout)
		if err == nil {
			if *conf.verbose {
				fmt.Printf("[%4d] Received: %v\n", i, ips)
			}

			if (len(ips) == 0 && *conf.acceptempty) || (len(ips) > 0) {
				totalAttempts = i
				break
			}
		} else if *conf.verbose {
			fmt.Printf("[%4d] Lookup failed: %v (%s)\n", i, err, *conf.host)
		}

		time.Sleep(time.Duration(*conf.pause) * time.Millisecond)
	}
	if !*conf.silent {
		fmt.Printf("Waited: %dms (%d attempts, host: %s)\n", (time.Now().UnixNano()-runTime)/int64(time.Millisecond), totalAttempts, *conf.host)
	}

}

func lookupWithTimeout(host *string, timeout *int) ([]net.IPAddr, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(*timeout)*time.Millisecond)
	defer cancel()
	var r net.Resolver
	return r.LookupIPAddr(ctx, *host)
}
