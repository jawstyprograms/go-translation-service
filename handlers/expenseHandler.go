package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "expense-tracker/config"
    "expense-tracker/models"
    "github.com/jackc/pgx/v5"
)

// CreateExpense creates a new expense
func CreateExpense(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    conn, err := config.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    var expense models.Expense
    err = json.NewDecoder(r.Body).Decode(&expense)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    query := `INSERT INTO expenses (description, amount, category, date) VALUES ($1, $2, $3, $4) RETURNING id`
    err = conn.QueryRow(context.Background(), query, expense.Description, expense.Amount, expense.Category, expense.Date).Scan(&expense.ID)
    if err != nil {
        http.Error(w, "Failed to create expense", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(expense)
}

// GetExpense retrieves an expense by ID
func GetExpense(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    conn, err := config.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid expense ID", http.StatusBadRequest)
        return
    }

    query := `SELECT id, description, amount, category, date FROM expenses WHERE id = $1`
    var expense models.Expense
    err = conn.QueryRow(context.Background(), query, id).Scan(
        &expense.ID,
        &expense.Description,
        &expense.Amount,
        &expense.Category,
        &expense.Date,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            http.Error(w, "Expense not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to get expense", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(expense)
}

// UpdateExpense updates an existing expense
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    conn, err := config.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid expense ID", http.StatusBadRequest)
        return
    }

    var expense models.Expense
    err = json.NewDecoder(r.Body).Decode(&expense)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    query := `UPDATE expenses SET description = $1, amount = $2, category = $3, date = $4 WHERE id = $5`
    _, err = conn.Exec(context.Background(), query, expense.Description, expense.Amount, expense.Category, expense.Date, id)
    if err != nil {
        http.Error(w, "Failed to update expense", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// DeleteExpense deletes an expense by ID
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    conn, err := config.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid expense ID", http.StatusBadRequest)
        return
    }

    query := `DELETE FROM expenses WHERE id = $1`
    _, err = conn.Exec(context.Background(), query, id)
    if err != nil {
        http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// GetAllExpenses retrieves all expenses
func GetAllExpenses(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    conn, err := config.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    query := `SELECT id, description, amount, category, date FROM expenses`
    rows, err := conn.Query(context.Background(), query)
    if err != nil {
        http.Error(w, "Failed to get expenses", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var expenses []models.Expense
    for rows.Next() {
        var expense models.Expense
        err = rows.Scan(
            &expense.ID,
            &expense.Description,
            &expense.Amount,
            &expense.Category,
            &expense.Date,
        )
        if err != nil {
            http.Error(w, "Failed to scan expense", http.StatusInternalServerError)
            return
        }
        expenses = append(expenses, expense)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(expenses)
}