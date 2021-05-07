package triangulate

func worker(jobChan <-chan Image) {
	for job := range jobChan {
		callGenerator(job)
	}
}
