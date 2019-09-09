package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jamiealquiza/envy"
	"net"
	"os"
	"time"
)

type config struct {
	host        *string
	timeout     *int
	retries     *int
	success     *int
	pause       *int64
	verbose     *bool
	silent      *bool
	acceptempty *bool
	exitcode    *bool
}

func main() {

	gracefulStop()

	runTime := time.Now().UnixNano()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nMIT License. Copyright (c) 2019 JamesJJ. https://github.com/JamesJJ/dns-ready\n")
	}

	conf := config{
		flag.String("host", "kube-dns.kube-system.svc.cluster.local", "hostname to resolve"),
		flag.Int("timeout", 1200, "Timeout in milliseconds for each lookup"),
		flag.Int("retries", 30, "Maximum number of attempts to resolve hostname"),
		flag.Int("success", 2, "Minimum number of lookups to count as success"),
		flag.Int64("pause", 800, "Number of milliseconds to pause between attempts"),
		flag.Bool("verbose", false, "Show process steps"),
		flag.Bool("silent", false, "Do not print anything to the console"),
		flag.Bool("acceptempty", false, "Accept a DNS response with no IP addresses"),
		flag.Bool("exitcode", false, "Exit with non-zero status if unsuccessful"),
	}
	envy.Parse("DNSREADY")
	flag.Parse()

	successCount := 0
	totalAttempts := *conf.retries
	for i := 1; i <= *conf.retries; i++ {
		ips, err := lookupWithTimeout(conf.host, conf.timeout)
		if err == nil {
			if *conf.verbose {
				fmt.Printf("[%4d] Received: %v\n", i, ips)
			}

			if (len(ips) == 0 && *conf.acceptempty) || (len(ips) > 0) {
				totalAttempts = i
				successCount++
				if successCount >= *conf.success {
					break
				}
			}
		} else if *conf.verbose {
			fmt.Printf("[%4d] Lookup failed: %v (%s)\n", i, err, *conf.host)
		}

		time.Sleep(time.Duration(*conf.pause) * time.Millisecond)
	}
	if !*conf.silent {
		fmt.Printf("Waited: %dms (success: %d / %d, host: %s)\n", (time.Now().UnixNano()-runTime)/int64(time.Millisecond), successCount, totalAttempts, *conf.host)
	}
	if (successCount < *conf.success) && (*conf.exitcode) {
		os.Exit(1)
	}

	os.Exit(0)

}

func lookupWithTimeout(host *string, timeout *int) ([]net.IPAddr, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(*timeout)*time.Millisecond)
	defer cancel()
	var r net.Resolver
	return r.LookupIPAddr(ctx, *host)
}
