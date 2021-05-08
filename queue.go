package triangulate

import "log"

func initQueue() {
	for i := 0; i < workerCount; i++ {
		log.Printf("Started queue worker %d\n", i)
		go worker(jobChan)
	}
	for i := 0; i < premiumWorkerCount; i++ {
		log.Printf("Started premium queue worker %d\n", i)
		go worker(premiumJobChan)
	}

}

func worker(jobChan <-chan Image) {
	for job := range jobChan {
		callGenerator(job)
	}
}
