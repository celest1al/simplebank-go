package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/celest1al/simplebank-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if valid := server.validAccount(ctx, req.FromAccountID, req.Currency); !valid {
		return
	}

	if valid := server.validAccount(ctx, req.ToAccountID, req.Currency); !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transfer successfully created",
		"data":    result,
	})
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account ID %d mismatch currency: %s vs %s", account.ID, account.Currency, currency)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
