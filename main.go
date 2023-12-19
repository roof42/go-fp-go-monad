package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/IBM/fp-go"
)

type Either struct {
	isRight bool
	value   interface{}
}

func Right(value interface{}) Either {
	return Either{isRight: true, value: value}
}

func Left(value interface{}) Either {
	return Either{isRight: false, value: value}
}

func readCSVFile(filePath string) Either {
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return Left(fmt.Sprintf("Error: Unable to open file - %v", err))
		}
		defer file.Close()

		reader := csv.NewReader(file)
		data, err := reader.ReadAll()
		if err != nil {
			return Left(fmt.Sprintf("Error: Unable to read CSV file - %v", err))
		}

		return Right(data)
	} else {
		return Left("Error: File not found")
	}
}

func removeRow(rowIndex int, data [][]string) Either {
	if len(data) > 1 {
		return Right(data[rowIndex:])
	} else {
		return Left("Error: Unable to remove header")
	}
}

func extractColumn(columnIndex int, data [][]string) Either {
	if len(data) > columnIndex {
		var columnValues []string
		for _, row := range data {
			columnValues = append(columnValues, row[columnIndex])
		}
		return Right(columnValues)
	} else {
		return Left("Error: Unable to extract column")
	}
}

func convertToFloat(data []string) Either {
	var convertedData []float64
	for _, item := range data {
		if floatValue, err := strconv.ParseFloat(item, 64); err == nil {
			convertedData = append(convertedData, floatValue)
		} else {
			return Left("Error: Unable to convert to float")
		}
	}
	return Right(convertedData)
}

func calculateAverage(columnValues []float64) Either {
	if len(columnValues) > 0 {
		sum := 0.0
		for _, value := range columnValues {
			sum += value
		}
		average := sum / float64(len(columnValues))
		return Right(average)
	} else {
		return Left("Error: Division by zero")
	}
}

func main() {
	csvFilePath := "example.csv"
	scoreColumnIndex := 1
	headerRowIndex := 1

	// Function composition using fp-go Pipe2
	result := fp.Pipe2(
		readCSVFile(csvFilePath),
		extractColumn(scoreColumnIndex),
		removeRow(headerRowIndex),
		convertToFloat,
		calculateAverage,
	)

	// Final result with if and else
	if result.(Either).isRight {
		fmt.Printf("An average score is %v\n", result.(Either).value)
	} else {
		fmt.Printf("Error processing data: %v\n", result.(Either).value)
	}
}
