package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Book struct {
	ID         string `json:"id"`
	VolumeInfo struct {
		Title       string   `json:"title"`
		Authors     []string `json:"authors"`
		Description string   `json:"description"`
		ImageLinks  struct {
			Thumbnail string `json:"thumbnail"`
		} `json:"imageLinks"`
	} `json:"volumeInfo"`
	SaleInfo struct {
		ListPrice struct {
			Amount float64 `json:"amount"`
		} `json:"listPrice"`
	} `json:"saleInfo"`
}

type SearchResult struct {
	Items []Book `json:"items"`
}

type PageData struct {
	BestSellers  []Book
	TrendingNow  []Book
	AwardWinners []Book
}

func fetchBooks(category string) ([]Book, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s", category))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []Book `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// fetch books from API based on category
func fetchBooksByCategory(category string) ([]Book, error) {
	apiURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=subject:%s&maxResults=20", url.QueryEscape(category))
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// cart data
func cartDataHandler(w http.ResponseWriter, r *http.Request) {
	cartItems := []Book{
		{
			ID: "wF6qDgAAQBAJ",
			VolumeInfo: struct {
				Title       string   `json:"title"`
				Authors     []string `json:"authors"`
				Description string   `json:"description"`
				ImageLinks  struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
			}{
				Title: "To Kill a Mockingbird",
				Authors: []string{
					"Harper Lee",
				},
				Description: "To Kill a Mockingbird is a novel by Harper Lee published in 1960. It was immediately successful, winning the Pulitzer Prize, and has become a classic of modern American literature.",
				ImageLinks: struct {
					Thumbnail string `json:"thumbnail"`
				}{
					Thumbnail: "/images/mockin.jpg",
				},
			},
			SaleInfo: struct {
				ListPrice struct {
					Amount float64 `json:"amount"`
				} `json:"listPrice"`
			}{
				ListPrice: struct {
					Amount float64 `json:"amount"`
				}{
					Amount: 7.99,
				},
			},
		},
		{
			ID: "RjsTAAAAYAAJ",
			VolumeInfo: struct {
				Title       string   `json:"title"`
				Authors     []string `json:"authors"`
				Description string   `json:"description"`
				ImageLinks  struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
			}{
				Title: "1984",
				Authors: []string{
					"George Orwell",
				},
				Description: "Nineteen Eighty-Four, often referred to as 1984, is a dystopian social science fiction novel by English novelist George Orwell. It was published on 8 June 1949 by Secker & Warburg as Orwell's ninth and final book completed in his lifetime.",
				ImageLinks: struct {
					Thumbnail string `json:"thumbnail"`
				}{
					Thumbnail: "/images/1984.png",
				},
			},
			SaleInfo: struct {
				ListPrice struct {
					Amount float64 `json:"amount"`
				} `json:"listPrice"`
			}{
				ListPrice: struct {
					Amount float64 `json:"amount"`
				}{
					Amount: 6.99,
				},
			},
		},
	}

	jsonData, err := json.Marshal(cartItems)
	if err != nil {
		http.Error(w, "Failed to fetch cart data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.Handle("/header/", http.StripPrefix("/header/", http.FileServer(http.Dir("header"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/forms/", http.StripPrefix("/forms/", http.FileServer(http.Dir("forms"))))
	http.Handle("/types/", http.StripPrefix("/types/", http.FileServer(http.Dir("types"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		bestSellers, _ := fetchBooks("bestsellers")
		trendingNow, _ := fetchBooks("trending")
		awardWinners, _ := fetchBooks("awardwinners")

		data := PageData{
			BestSellers:  bestSellers,
			TrendingNow:  trendingNow,
			AwardWinners: awardWinners,
		}

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/cart/add", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.URL.Query().Get("bookId")
		fmt.Printf("Book added to cart: %s\n", bookId)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("header/about.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("header/contact.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/faq", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("header/faq.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("header/cart.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("forms/login.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("forms/signup.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		apiUrl := "https://www.googleapis.com/books/v1/volumes?q=" + url.QueryEscape(query)
		resp, err := http.Get(apiUrl)
		if err != nil {
			http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var result SearchResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, "Failed to decode response", http.StatusInternalServerError)
			return
		}

		// Загружаем и парсим шаблон search.html
		tmpl, err := template.ParseFiles("header/search.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, result.Items); err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/fiction", func(w http.ResponseWriter, r *http.Request) {
		books, err := fetchBooksByCategory("fiction")
		if err != nil {
			http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("types/fiction.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, books)
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/nonfiction", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("types/nonfiction.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/science", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("types/science.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/fantasy", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("types/fantasy.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/mystery", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("types/mystery.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
