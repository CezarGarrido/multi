package multi_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/CezarGarrido/multi"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	portDB   = "5432"
	user     = "postgres"
	password = "postgres"
	dbname   = "locadora"
)

func TestMultiRun(t *testing.T) {

	db := NewDatabase()

	ctx := context.Background()

	values, err := multi.New().Run("Inserir locadora e cliente", db, func(tx *sql.Tx) (multi.Result, error) {

		locadora := NewLocadoraWithTx(tx)

		locadoraID, err := locadora.CreateWithTx(ctx, "Joao", "joa@joao", "55 9999-9999")
		if err != nil {
			return nil, err
		}

		cliente := NewClienteWithTx(tx)

		clienteID, err := cliente.CreateWithTx(ctx, locadoraID, "Maria", "01234567890")
		if err != nil {
			return nil, err
		}

		result := multi.Result{
			"locadora_id": locadoraID,
			"cliente_id":  clienteID,
		}

		return result, nil
	})

	if err != nil {
		t.Error("Erro ao adicionar locadora e cliente", err.Error())
		return
	}

	fmt.Println(values)
}

func NewDatabase() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, portDB, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

type Locadora struct {
	Conn   *sql.DB
	ConnTx *sql.Tx
}

func NewLocadoraWithTx(connTx *sql.Tx) *Locadora {
	return &Locadora{
		ConnTx: connTx,
	}
}

// CreateWithTx : ...
func (locadora *Locadora) CreateWithTx(ctx context.Context, nome, email, telefone string) (int64, error) {

	query := `INSERT INTO locadoras (nome, email, telefone,  created_at) VALUES($1, $2, $3, $4) RETURNING id`
	stmt, err := locadora.ConnTx.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var ID int64
	err = stmt.QueryRowContext(ctx,
		nome,
		email,
		telefone,
		time.Now(),
	).Scan(&ID)

	if err != nil {
		return -1, err
	}

	return ID, nil
}

type Cliente struct {
	Conn   *sql.DB
	ConnTx *sql.Tx
}

func NewClienteWithTx(connTx *sql.Tx) *Cliente {
	return &Cliente{
		ConnTx: connTx,
	}
}

// CreateWithTx : ...
func (cliente *Cliente) CreateWithTx(ctx context.Context, locadoraID int64, nome, cpf string) (int64, error) {

	query := `INSERT INTO clientes (locadora_id, nome, cpf, created_at) VALUES($1, $2, $3, $4) RETURNING id`
	stmt, err := cliente.ConnTx.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var ID int64
	err = stmt.QueryRowContext(ctx,
		locadoraID,
		nome,
		cpf,
		time.Now(),
	).Scan(&ID)

	if err != nil {
		return -1, err
	}

	return ID, nil
}
