package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/muhammadzaid-99/SubSnip/internal/models"
	"github.com/muhammadzaid-99/SubSnip/internal/queue"
	"github.com/muhammadzaid-99/SubSnip/internal/status"
)

func init() {
	// queue.Init()
	status.Init()
}

func main() {
	ch := queue.GetChannel()

	msgs, err := ch.Consume(
		queue.GetQueueName(), "", false, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Worker shutting down...")
				return
			default:
				// No shutdown signal -> try to receive a message (blocking)
			}

			d, ok := <-msgs
			if !ok {
				log.Printf("✔️  Worker channel closed.")
				return
			}

			log.Printf("👷 Worker got task")

			var task models.TaskRequest
			if err := json.Unmarshal(d.Body, &task); err != nil {
				d.Nack(false, false)
				continue
			}

			// status.Set(task.TaskID, "processing")
			log.Println("Processing task: ", task.TaskID)
			taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			err := runPythonScript(taskCtx, &task) // pass the context
			cancel()                               // always cancel to release resources
			if err != nil {
				log.Println("Task failed.")
				// status.Set(task.TaskID, "failed")
			} else {
				log.Println("Task completed.")
				// status.Set(task.TaskID, "finished")
			}
			d.Ack(false)

		}
	}()
	<-ctx.Done()
	log.Println("🛑 Shutdown signal received. Waiting for worker to complete.")
	queue.Close()
}

func runPythonScript(ctx context.Context, task *models.TaskRequest) error {
	cmd := exec.CommandContext(ctx, "python", "scripts/extract_subs.py")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go func() {
		defer stdin.Close()
		json.NewEncoder(stdin).Encode(task)
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Printf("[PYTHON OUT] %s", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("[PYTHON ERR] %s", scanner.Text())
		}
	}()

	err := cmd.Start()
	if err != nil {
		log.Printf(" Failed to start Python script: %v", err)
		return err
	}

	err = cmd.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Task %s timed out", task.TaskID)
		return ctx.Err()
	}

	if err != nil {
		log.Printf("Python script failed: %v", err)
		return err
	}

	log.Printf("Python script finished successfully")
	return nil
}

// func StartWorkers(doneWorkers chan struct{}) {
// 	ch := queue.GetChannel()

// 	msgs, err := ch.Consume(
// 		queue.GetQueueName(), "", false, false, false, false, nil,
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register consumer: %v", err)
// 	}

// 	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	const maxWorkers = 5
// 	var wg sync.WaitGroup

// 	for i := range maxWorkers {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					log.Printf("Worker %d shutting down...", workerID)
// 					return
// 				default:
// 					// No shutdown signal -> try to receive a message (blocking)
// 				}

// 				d, ok := <-msgs
// 				if !ok {
// 					log.Printf(" Worker %d: message channel closed", workerID)
// 					return
// 				}

// 				log.Printf("Worker %d got task", workerID)

// 				var task models.TaskRequest
// 				if err := json.Unmarshal(d.Body, &task); err != nil {
// 					d.Nack(false, false)
// 					continue
// 				}

// 				status.Set(task.TaskID, "processing")
// 				log.Println("Processing task: ", task.TaskID)
// 				taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
// 				err := runPythonScript(taskCtx, &task) // pass the context
// 				cancel()                               // always cancel to release resources
// 				if err != nil {
// 					status.Set(task.TaskID, "failed")
// 				} else {
// 					status.Set(task.TaskID, "finished")
// 				}
// 				d.Ack(false)

// 			}
// 		}(i + 1)
// 	}

// 	log.Println("Consumer started...")

// 	<-ctx.Done()
// 	log.Println("Shutdown signal received. Waiting for workers...")

// 	// Wait for all workers to finish
// 	wg.Wait()
// 	log.Println("All workers exited.")
// 	doneWorkers <- struct{}{}
// }
