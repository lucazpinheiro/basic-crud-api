package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "crud-db"
)

type Product struct {
	ID          int     `json:",omitempty"`
	Name        string  `json:",omitempty"`
	Available   bool    `json:",omitempty"`
	Price       float64 `json:",omitempty"`
	Description string  `json:",omitempty"`
}

func main() {
	db, err := OpenConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/products", func(c *fiber.Ctx) error {
		var product Product
		var products []Product
		sql := `SELECT id, name, isAvailable, price, description FROM public.products`
		rows, err := db.Query(sql)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&product.ID, &product.Name, &product.Available, &product.Price, &product.Description)
			if err != nil {
				log.Fatal(err)
				c.JSON("An error occured")
			}
			products = append(products, product)
		}

		return c.JSON(products)
	})

	app.Get("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var product Product
		sql := `SELECT id, name, isAvailable, price, description FROM public.products WHERE id = $1`
		row := db.QueryRow(sql, id)

		err := row.Scan(&product.ID, &product.Name, &product.Available, &product.Price, &product.Description)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}

		return c.JSON(product)
	})

	app.Post("/products", func(c *fiber.Ctx) error {
		p := new(Product)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		sql := `INSERT INTO public.products (name, isAvailable, price, description) VALUES ($1, $2, $3, $4) RETURNING id`
		err := db.QueryRow(sql, p.Name, p.Available, p.Price, p.Description).Scan(&p.ID)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}

		return c.JSON(&p.ID)
	})

	app.Put("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var p Product

		var productID int
		selectSQL := `SELECT id FROM public.products WHERE id = $1`
		err := db.QueryRow(selectSQL, id).Scan(&productID)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}
		if err := c.BodyParser(&p); err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}

		updateSQL := `UPDATE public.products SET name = $1, isAvailable = $2, price = $3, description = $4 WHERE id = $5`
		_, err = db.Exec(updateSQL, p.Name, p.Available, p.Price, p.Description, id)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}

		return c.JSON(p)
	})

	app.Delete("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var deletedID int

		sql := `DELETE FROM public.products WHERE id = $1`
		err := db.QueryRow(sql, id).Scan(&deletedID)
		if err != nil {
			log.Fatal(err)
			c.JSON("An error occured")
		}

		return c.JSON(deletedID)
	})

	app.Listen(":3000")
}

func OpenConn() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	return db, err
}
