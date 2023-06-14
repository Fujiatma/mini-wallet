package controllers

import (
	"errors"
	"github.com/julo/mini-wallet/src/middleware"
	"github.com/julo/mini-wallet/src/models"
	resp "github.com/julo/mini-wallet/src/response"
	"github.com/julo/mini-wallet/src/services"
	"net/http"
	"strconv"
)

type WalletController struct {
	WalletService   *services.WalletService
	CustomerService *services.AuthService
}

func NewWalletController(walletService *services.WalletService, customerService *services.AuthService) *WalletController {
	return &WalletController{
		WalletService:   walletService,
		CustomerService: customerService,
	}
}

func (c *WalletController) getJWTClaims(r *http.Request) (*middleware.JWTClaims, error) {
	claims, ok := r.Context().Value(middleware.JWTClaimsContextKey).(*middleware.JWTClaims)
	if !ok || claims == nil {
		return nil, errors.New("Invalid JWT claims")
	}
	return claims, nil
}

func (c *WalletController) EnableWalletController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	customer, err := c.CustomerService.GetCustomerByCustomerXID(claims.CustomerXID)
	if err != nil {
		http.Error(w, "Failed to get customer", http.StatusInternalServerError)
		return
	}

	// Aktifkan wallet jika belum diaktifkan
	wallet, err := c.WalletService.EnableWallet(customer.CustomerXID)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyEnabledError); ok && customErr.Code == 800 {
			http.Error(w, "Wallet is already enabled", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to enable wallet", http.StatusInternalServerError)
		return
	}

	// Mengonstruksi response
	response := map[string]interface{}{
		"wallet": wallet,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)
}

func (c *WalletController) GetWalletBalanceController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Mendapatkan informasi saldo dompet
	customer, err := c.CustomerService.GetCustomerByCustomerXID(claims.CustomerXID)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyDisabledError); ok && customErr.Code == 801 {
			http.Error(w, "Wallet is already disabled", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to get customer wallet", http.StatusInternalServerError)
		return
	}

	if customer == nil {
		http.Error(w, "Failed to get customer wallet", http.StatusNotFound)
		return
	}

	// Membangun respons
	response := map[string]interface{}{
		"wallet": customer.Wallet,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)

}

func (c *WalletController) GetWalletTransactionController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Mengambil transaksi wallet dari repository berdasarkan customerID
	transactions, err := c.WalletService.GetWalletTransactions(claims.CustomerXID)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyDisabledError); ok && customErr.Code == 801 {
			http.Error(w, "Wallet is already disabled", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to get wallet transactions", http.StatusInternalServerError)
		return
	}

	// Mengirimkan respons JSON dengan daftar transaksi wallet
	response := struct {
		Status string               `json:"status"`
		Data   []models.Transaction `json:"data"`
	}{
		Status: "success",
		Data:   transactions,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)
}

func (c *WalletController) AddVirtualMoneyController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Membaca data deposit dari permintaan
	amountStr := r.FormValue("amount")
	referenceID := r.FormValue("reference_id")

	// Memvalidasi jumlah deposit
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	if amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Melakukan penambahan virtual money ke wallet
	deposit, err := c.WalletService.AddVirtualMoney(claims.CustomerXID, amount, referenceID)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyDisabledError); ok && customErr.Code == 801 {
			http.Error(w, "Wallet is already disabled", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to add virtual money to wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mengirimkan respons JSON dengan data deposit
	response := struct {
		Status string          `json:"status"`
		Data   *models.Deposit `json:"data"`
	}{
		Status: "success",
		Data:   deposit,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)
}

func (c *WalletController) UseVirtualMoneyController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Membaca data deposit dari permintaan
	amountStr := r.FormValue("amount")
	referenceID := r.FormValue("reference_id")

	// Memvalidasi jumlah deposit
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	if amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Melakukan penggunaan virtual money dari wallet
	withdrawal, err := c.WalletService.UseVirtualMoney(claims.CustomerXID, amount, referenceID)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyDisabledError); ok && customErr.Code == 801 {
			http.Error(w, "Wallet is already disabled", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to use virtual money from wallet", http.StatusInternalServerError)
		return
	}

	// Mengirimkan respons JSON dengan data penarikan
	response := struct {
		Status string             `json:"status"`
		Data   *models.Withdrawal `json:"data"`
	}{
		Status: "success",
		Data:   withdrawal,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)
}

func (c *WalletController) DisableWalletController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan klaim JWT dari konteks permintaan
	claims, err := c.getJWTClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Mendapatkan data disable wallet dari body permintaan
	isDisabled := r.FormValue("is_disabled")

	// Melakukan validasi data disable wallet
	isDisabledBool, err := strconv.ParseBool(isDisabled)
	if err != nil {
		http.Error(w, "Invalid value for is_disabled", http.StatusBadRequest)
		return
	}

	// Menonaktifkan wallet
	wallet, err := c.WalletService.DisableWallet(claims.CustomerXID, isDisabledBool)
	if err != nil {
		if customErr, ok := err.(resp.WalletAlreadyDisabledError); ok && customErr.Code == 801 {
			http.Error(w, "Wallet is already disabled", http.StatusOK)
			return
		}
		http.Error(w, "Failed to disable wallet", http.StatusInternalServerError)
		return
	}

	// Mengirimkan respons JSON dengan data wallet
	response := struct {
		Status string         `json:"status"`
		Data   *models.Wallet `json:"data"`
	}{
		Status: "success",
		Data:   wallet,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)

}
