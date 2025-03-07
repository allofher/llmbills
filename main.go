package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Bill represents the legislative bill data
type Bill struct {
	ID           string `xml:"id"`
	Title        string `xml:"title"`
	ShortTitle   string `xml:"shortTitle"`
	Parliament   string `xml:"parliament"`
	Session      string `xml:"session"`
	DateRange    string `xml:"dateRange"`
	Sponsor      string `xml:"sponsor"`
	BillType     string `xml:"billType"`
	CurrentStatus string `xml:"currentStatus"`
	Content      string `xml:"content"`
}

// PageData holds all data needed for page rendering
type PageData struct {
	Bill Bill
	XMLContent string
}

func main() {
	// Create file server for static assets
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/bill/", billHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/donate", donateHandler)


	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// In a real app, you would fetch this data from a database
	bill := Bill{
		ID:           "S-204",
		Title:        "An Act to amend the Customs Tariff (goods from Xinjiang)",
		ShortTitle:   "Xinjiang Manufactured Goods Importation Prohibition Act",
		Parliament:   "44th Parliament",
		Session:      "1st session",
		DateRange:    "November 22, 2023 to January 6, 2025",
		Sponsor:      "Sen. Leo Housakos",
		BillType:     "Senate Public Bill",
		CurrentStatus: "At second reading in the Senate",
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", bill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func billHandler(w http.ResponseWriter, r *http.Request) {
	// Extract bill ID from URL
	billID := r.URL.Path[len("/bill/"):]
	if billID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// In a real app, you would fetch the bill from a database
	// For this example, we'll use a hardcoded bill if it matches S-204
	if billID != "S-204" {
		http.NotFound(w, r)
		return
	}

	// Load XML content
	xmlContent, err := os.ReadFile("data/S-204.xml")
	if err != nil {
		http.Error(w, "Failed to load bill XML", http.StatusInternalServerError)
		return
	}

	var bill Bill
	err = xml.Unmarshal(xmlContent, &bill)
	if err != nil {
		http.Error(w, "Failed to parse bill XML", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Bill: bill,
		XMLContent: string(xmlContent),
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/bill.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	// Extract bill ID from URL
	billID := r.URL.Path[len("/bill/pdf/"):]
	if billID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// In a real app, you would generate a PDF from the XML
	// For this example, we'll serve a static PDF file
	pdfPath := filepath.Join("data", billID+".pdf")
	
	// Check if the file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set content type and headers for PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.pdf\"", billID))

	// Serve the PDF file
	http.ServeFile(w, r, pdfPath)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// In a real app, you would save the email to a database
	// For this example, we'll just log it
	log.Printf("New subscription: %s", email)

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	searchterm := r.FormValue("search term")
	if searchterm == "" {
		http.Error(w, "search term required", http.StatusBadRequest)
		return
	}

	// In a real app, you would save the email to a database
	// For this example, we'll just log it
	log.Printf("New search: %s", searchterm)

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func donateHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, you would redirect to a payment processor
	// For this example, we'll redirect to Google as requested
	http.Redirect(w, r, "https://www.google.com", http.StatusSeeOther)
}
