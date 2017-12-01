package main

// Note: Most Go-plugins for editors maintain the import list automatically.
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	rows := readOrders("orders.csv")
	rows = calculate(rows)
	writeOrders("ordersReport.csv", rows)
}

func readOrders(name string) [][]string {

	f, err := os.Open(name)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	rows, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}

	return rows
}

func calculate(rows [][]string) [][]string {

	sum := 0
	nb := 0

	for i := range rows {

		if i == 0 { // the first row
			rows[0] = append(rows[0], "Total")
			continue // start the next loop iteration
		}

		// All subsequent rows:

		item := rows[i][2] // Column 2 contains the item name

		// Column 3 contains the price. Remove the decimal point, and
		// turn the value into an integer (representing the value in cents)
		price, err := strconv.Atoi(strings.Replace(rows[i][3], ".", "", -1))
		if err != nil {
			log.Fatalf("Cannot retrieve price of %s: %s\n", item, err)
		}

		// Column 4 contains the ordered quantity.
		qty, err := strconv.Atoi(rows[i][4])
		if err != nil {
			log.Fatalf("Cannot retrieve quantity of %s: %s\n", item, err)
		}

		// Calculate the total and append it to the current row.
		total := price * qty
		rows[i] = append(rows[i], intToFloatString(total))
		sum += total

		// Update the total sum and the # of ball pens.
		if item == "Ball Pen" {
			nb += qty
		}
	}

	// Here we create two new rows. The first one shows the total sum, and
	// the second one shows the number of ball pens ordered.
	rows = append(rows, []string{"", "", "", "Sum", "", intToFloatString(sum)})
	rows = append(rows, []string{"", "", "", "Ball Pens", fmt.Sprint(nb), ""})

	return rows
}

func intToFloatString(n int) string {
	intgr := n / 100
	frac := n - intgr*100
	return fmt.Sprintf("%d.%d", intgr, frac)
}

func writeOrders(name string, rows [][]string) {

	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	// We are going to write to a file, so any errors are relevant and
	// need to be logged. Hence the anonymous func instead of a one-liner.
	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Cannot close '%s': %s\n", name, e.Error())
		}
	}()

	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
}
