package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/epicira/go-iracluster"
)

func main() {
	iraCluster := iracluster.New("app_name", os.Getenv("IRA_CLUSTER_KEY"), func(msg string) {
		fmt.Println("event:")
		fmt.Println(msg)
	})
	iraCluster.JoinCluster()

	test_db := iracluster.NewTDB(iraCluster, "test_db", "CREATE TABLE info(x INT, y INT);", "", true, nil)
	status := test_db.Open()
	fmt.Printf("\nTDB Open Status: %s\n", status)

	q1 := "INSERT INTO info (x, y) VALUES (1, 2)"
	fmt.Println(test_db.Execute(q1))
	fmt.Printf("Executed q1: %s\n", q1)

	q2 := "INSERT INTO info (x, y) VALUES (3, 4)"
	delay := 3000 //ms
	test_db.ExecuteAsync(q2, delay)
	fmt.Printf("Executed q2: %s\n", q2)

	fmt.Printf("Row Count: %d\n", test_db.Count("info", ""))
	fmt.Println("sleeping for 5 seconds...")
	time.Sleep(5 * time.Second)
	fmt.Printf("Row Count: %d\n", test_db.Count("info", ""))

	fmt.Println("All rows:")
	fmt.Println(test_db.Select("SELECT * FROM info"))

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
