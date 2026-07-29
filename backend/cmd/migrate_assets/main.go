package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Migrate Gold
	rows, err := pool.Query(ctx, "SELECT id, name FROM assets WHERE category = 'gold' AND (ticker IS NULL OR ticker = '')")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var goldUpdates []map[string]interface{}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		parts := strings.Split(name, " - ")
		if len(parts) >= 2 {
			brand := strings.TrimSpace(parts[0])
			rest := strings.TrimSpace(strings.Join(parts[1:], " - "))
			
			// Remove the unit part like (1 Chỉ)
			unitRegex := regexp.MustCompile(`\(.+?\)$`)
			purityStr := strings.TrimSpace(unitRegex.ReplaceAllString(rest, ""))
			
			typeCode := "NHAN"
			if strings.Contains(purityStr, "Vàng miếng") {
				typeCode = "MIENG"
			} else if strings.Contains(purityStr, "Vàng trang sức") {
				typeCode = "TS"
			}

			brandCode := "KHAC"
			if brand == "Bảo Tín Minh Châu" {
				brandCode = "BTMC"
			} else if brand == "Bảo Tín Mạnh Hải" {
				brandCode = "BTMH"
			} else if brand == "Phú Quý" {
				brandCode = "PHUQUY"
			} else if brand == "Mi Hồng" {
				brandCode = "MIHONG"
			} else if brand != "Tư nhân" && brand != "Khác" {
				brandCode = strings.ReplaceAll(strings.ToUpper(brand), " ", "")
			}

			ticker := fmt.Sprintf("GOLD_%s_%s", brandCode, typeCode)
			if len(ticker) > 20 {
				ticker = ticker[:20]
			}
			goldUpdates = append(goldUpdates, map[string]interface{}{"id": id, "ticker": ticker})
		}
	}

	for _, u := range goldUpdates {
		_, err := pool.Exec(ctx, "UPDATE assets SET ticker = $1 WHERE id = $2", u["ticker"], u["id"])
		if err != nil {
			log.Printf("Failed to update asset %s: %v", u["id"], err)
		} else {
			fmt.Printf("Updated Gold %s with ticker %s\n", u["id"], u["ticker"])
		}
	}

	// Migrate Fund
	rows2, err := pool.Query(ctx, "SELECT id, name FROM assets WHERE category = 'fund' AND (ticker IS NULL OR ticker = '')")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	var fundUpdates []map[string]interface{}
	for rows2.Next() {
		var id, name string
		if err := rows2.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		parts := strings.Split(name, " - ")
		if len(parts) >= 1 {
			ticker := strings.TrimSpace(parts[0])
			fundUpdates = append(fundUpdates, map[string]interface{}{"id": id, "ticker": ticker})
		}
	}

	for _, u := range fundUpdates {
		_, err := pool.Exec(ctx, "UPDATE assets SET ticker = $1 WHERE id = $2", u["ticker"], u["id"])
		if err != nil {
			log.Printf("Failed to update asset %s: %v", u["id"], err)
		} else {
			fmt.Printf("Updated Fund %s with ticker %s\n", u["id"], u["ticker"])
		}
	}

	fmt.Println("Migration complete!")
}
