package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type BenchmarkResult struct {
	TotalRequests   int
	SuccessRequests int
	FailedRequests  int
	Duration        time.Duration
	RequestsPerSec  float64
	AvgLatency      time.Duration
	MinLatency      time.Duration
	MaxLatency      time.Duration
}

func main() {
	baseURL := "http://localhost:8080"
	concurrency := 10
	totalRequests := 1000

	fmt.Printf("Benchmarking Todo Microservice\n")
	fmt.Printf("URL: %s\n", baseURL)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Total Requests: %d\n\n", totalRequests)

	fmt.Println("=== Signup Benchmark ===")
	signupResult := benchmarkSignup(baseURL, concurrency, totalRequests)
	printResults(signupResult)

	fmt.Println("\n=== Login Benchmark ===")
	loginResult := benchmarkLogin(baseURL, concurrency, totalRequests)
	printResults(loginResult)

	fmt.Println("\n=== Create Todo Benchmark ===")
	createTodoResult := benchmarkCreateTodo(baseURL, concurrency, totalRequests)
	printResults(createTodoResult)

	fmt.Println("\n=== List Todos Benchmark ===")
	listTodosResult := benchmarkListTodos(baseURL, concurrency, totalRequests)
	printResults(listTodosResult)
}

func benchmarkSignup(baseURL string, concurrency, totalRequests int) BenchmarkResult {
	return runBenchmark(baseURL+"/signup", "POST", func(i int) []byte {
		data := map[string]string{
			"email":    fmt.Sprintf("user%d@example.com", i),
			"password": "password123",
		}
		body, _ := json.Marshal(data)
		return body
	}, concurrency, totalRequests)
}

func benchmarkLogin(baseURL string, concurrency, totalRequests int) BenchmarkResult {
	http.Post(baseURL+"/signup", "application/json", bytes.NewBuffer([]byte(`{"email":"bench@example.com","password":"pass123"}`)))

	return runBenchmark(baseURL+"/login", "POST", func(i int) []byte {
		data := map[string]string{
			"email":    "bench@example.com",
			"password": "pass123",
		}
		body, _ := json.Marshal(data)
		return body
	}, concurrency, totalRequests)
}

func benchmarkCreateTodo(baseURL string, concurrency, totalRequests int) BenchmarkResult {
	return runBenchmark(baseURL+"/todos", "POST", func(i int) []byte {
		data := map[string]string{
			"user_id": "user_1",
			"text":    fmt.Sprintf("Todo item %d", i),
		}
		body, _ := json.Marshal(data)
		return body
	}, concurrency, totalRequests)
}

func benchmarkListTodos(baseURL string, concurrency, totalRequests int) BenchmarkResult {
	return runBenchmark(baseURL+"/todos?user_id=user_1&limit=50&offset=0", "GET", func(i int) []byte {
		return nil
	}, concurrency, totalRequests)
}

func runBenchmark(url, method string, bodyGenerator func(int) []byte, concurrency, totalRequests int) BenchmarkResult {
	var wg sync.WaitGroup
	var mu sync.Mutex

	successCount := 0
	failedCount := 0
	var totalLatency time.Duration
	minLatency := time.Hour
	maxLatency := time.Duration(0)

	requestsPerWorker := totalRequests / concurrency
	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 10 * time.Second}

			for j := 0; j < requestsPerWorker; j++ {
				requestStart := time.Now()

				var req *http.Request
				var err error

				body := bodyGenerator(workerID*requestsPerWorker + j)
				if body != nil {
					req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(method, url, nil)
				}

				if err != nil {
					mu.Lock()
					failedCount++
					mu.Unlock()
					continue
				}

				resp, err := client.Do(req)
				latency := time.Since(requestStart)

				mu.Lock()
				totalLatency += latency
				if latency < minLatency {
					minLatency = latency
				}
				if latency > maxLatency {
					maxLatency = latency
				}

				if err != nil || resp.StatusCode >= 400 {
					failedCount++
				} else {
					successCount++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	return BenchmarkResult{
		TotalRequests:   totalRequests,
		SuccessRequests: successCount,
		FailedRequests:  failedCount,
		Duration:        duration,
		RequestsPerSec:  float64(successCount) / duration.Seconds(),
		AvgLatency:      totalLatency / time.Duration(successCount+failedCount),
		MinLatency:      minLatency,
		MaxLatency:      maxLatency,
	}
}

func printResults(result BenchmarkResult) {
	fmt.Printf("Total Requests:    %d\n", result.TotalRequests)
	fmt.Printf("Successful:        %d\n", result.SuccessRequests)
	fmt.Printf("Failed:            %d\n", result.FailedRequests)
	fmt.Printf("Duration:          %v\n", result.Duration)
	fmt.Printf("Requests/sec:      %.2f\n", result.RequestsPerSec)
	fmt.Printf("Avg Latency:       %v\n", result.AvgLatency)
	fmt.Printf("Min Latency:       %v\n", result.MinLatency)
	fmt.Printf("Max Latency:       %v\n", result.MaxLatency)
}
