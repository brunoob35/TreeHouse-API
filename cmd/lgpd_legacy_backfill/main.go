package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/config"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

type backfillStats struct {
	RowsSeen            int
	NamesProtected      int
	CPFProtected        int
	RGProtected         int
	LegacyCPFAnonymized int
	LegacyRGAnonymized  int
}

func main() {
	config.LoadEnv()

	db, err := persistency.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := ensureRequiredColumns(db); err != nil {
		log.Fatal(err)
	}

	usersStats, err := backfillUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	customersStats, err := backfillCustomers(db)
	if err != nil {
		log.Fatal(err)
	}

	studentsStats, err := backfillStudents(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("LGPD legacy backfill concluido.\nusuarios=%+v\nclientes=%+v\nalunos=%+v",
		usersStats,
		customersStats,
		studentsStats,
	)
}

func ensureRequiredColumns(db *sql.DB) error {
	required := []struct {
		table  string
		column string
	}{
		{"usuarios", "nome_protegido"},
		{"usuarios", "cpf_protegido"},
		{"usuarios", "rg_protegido"},
		{"usuarios", "cpf_hash"},
		{"usuarios", "rg_hash"},
		{"clientes", "nome_protegido"},
		{"clientes", "cpf_protegido"},
		{"clientes", "rg_protegido"},
		{"clientes", "cpf_hash"},
		{"clientes", "rg_hash"},
		{"alunos", "nome_protegido"},
	}

	for _, item := range required {
		ok, err := columnExists(db, item.table, item.column)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf(
				"coluna obrigatoria ausente: %s.%s. Rode antes o patch SQL de LGPD em TreeHouse-API/sql/patch_lgpd_dev_prod.sql",
				item.table,
				item.column,
			)
		}
	}

	return nil
}

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	var exists int
	err := db.QueryRow(`
		SELECT 1
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = ?
		LIMIT 1
	`, tableName, columnName).Scan(&exists)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func backfillUsers(db *sql.DB) (backfillStats, error) {
	rows, err := db.Query(`
		SELECT
			id,
			COALESCE(nome, ''),
			COALESCE(cpf, ''),
			COALESCE(rg, ''),
			COALESCE(nome_protegido, ''),
			COALESCE(cpf_protegido, ''),
			COALESCE(rg_protegido, '')
		FROM usuarios
	`)
	if err != nil {
		return backfillStats{}, err
	}
	defer rows.Close()

	stats := backfillStats{}

	for rows.Next() {
		var id uint64
		var nome string
		var cpf string
		var rg string
		var nomeProtegido string
		var cpfProtegido string
		var rgProtegido string

		if err := rows.Scan(&id, &nome, &cpf, &rg, &nomeProtegido, &cpfProtegido, &rgProtegido); err != nil {
			return stats, err
		}
		stats.RowsSeen++

		updates := []string{}
		args := []interface{}{}

		if strings.TrimSpace(nome) != "" && strings.TrimSpace(nomeProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(nome))
			if encErr != nil {
				return stats, fmt.Errorf("usuarios.id=%d nome: %w", id, encErr)
			}
			updates = append(updates, "nome_protegido = ?")
			args = append(args, encrypted)
			stats.NamesProtected++
		}

		if strings.TrimSpace(cpf) != "" && strings.TrimSpace(cpfProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(cpf))
			if encErr != nil {
				return stats, fmt.Errorf("usuarios.id=%d cpf: %w", id, encErr)
			}
			hash, hashErr := security.HashPII(strings.TrimSpace(cpf))
			if hashErr != nil {
				return stats, fmt.Errorf("usuarios.id=%d cpf hash: %w", id, hashErr)
			}
			updates = append(updates, "cpf_protegido = ?", "cpf_hash = ?", "cpf = NULL")
			args = append(args, encrypted, hash)
			stats.CPFProtected++
			stats.LegacyCPFAnonymized++
		}

		if strings.TrimSpace(rg) != "" && strings.TrimSpace(rgProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(rg))
			if encErr != nil {
				return stats, fmt.Errorf("usuarios.id=%d rg: %w", id, encErr)
			}
			hash, hashErr := security.HashPII(strings.TrimSpace(rg))
			if hashErr != nil {
				return stats, fmt.Errorf("usuarios.id=%d rg hash: %w", id, hashErr)
			}
			updates = append(updates, "rg_protegido = ?", "rg_hash = ?", "rg = NULL")
			args = append(args, encrypted, hash)
			stats.RGProtected++
			stats.LegacyRGAnonymized++
		}

		if len(updates) == 0 {
			continue
		}

		args = append(args, id)
		query := fmt.Sprintf("UPDATE usuarios SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?", strings.Join(updates, ", "))
		if _, err := db.Exec(query, args...); err != nil {
			return stats, err
		}
	}

	return stats, rows.Err()
}

func backfillCustomers(db *sql.DB) (backfillStats, error) {
	rows, err := db.Query(`
		SELECT
			id,
			COALESCE(nome, ''),
			COALESCE(cpf, ''),
			COALESCE(rg, ''),
			COALESCE(nome_protegido, ''),
			COALESCE(cpf_protegido, ''),
			COALESCE(rg_protegido, '')
		FROM clientes
	`)
	if err != nil {
		return backfillStats{}, err
	}
	defer rows.Close()

	stats := backfillStats{}

	for rows.Next() {
		var id uint64
		var nome string
		var cpf string
		var rg string
		var nomeProtegido string
		var cpfProtegido string
		var rgProtegido string

		if err := rows.Scan(&id, &nome, &cpf, &rg, &nomeProtegido, &cpfProtegido, &rgProtegido); err != nil {
			return stats, err
		}
		stats.RowsSeen++

		updates := []string{}
		args := []interface{}{}

		if strings.TrimSpace(nome) != "" && strings.TrimSpace(nomeProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(nome))
			if encErr != nil {
				return stats, fmt.Errorf("clientes.id=%d nome: %w", id, encErr)
			}
			updates = append(updates, "nome_protegido = ?")
			args = append(args, encrypted)
			stats.NamesProtected++
		}

		if strings.TrimSpace(cpf) != "" && strings.TrimSpace(cpfProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(cpf))
			if encErr != nil {
				return stats, fmt.Errorf("clientes.id=%d cpf: %w", id, encErr)
			}
			hash, hashErr := security.HashPII(strings.TrimSpace(cpf))
			if hashErr != nil {
				return stats, fmt.Errorf("clientes.id=%d cpf hash: %w", id, hashErr)
			}
			updates = append(updates, "cpf_protegido = ?", "cpf_hash = ?", "cpf = NULL")
			args = append(args, encrypted, hash)
			stats.CPFProtected++
			stats.LegacyCPFAnonymized++
		}

		if strings.TrimSpace(rg) != "" && strings.TrimSpace(rgProtegido) == "" {
			encrypted, encErr := security.EncryptPII(strings.TrimSpace(rg))
			if encErr != nil {
				return stats, fmt.Errorf("clientes.id=%d rg: %w", id, encErr)
			}
			hash, hashErr := security.HashPII(strings.TrimSpace(rg))
			if hashErr != nil {
				return stats, fmt.Errorf("clientes.id=%d rg hash: %w", id, hashErr)
			}
			updates = append(updates, "rg_protegido = ?", "rg_hash = ?", "rg = NULL")
			args = append(args, encrypted, hash)
			stats.RGProtected++
			stats.LegacyRGAnonymized++
		}

		if len(updates) == 0 {
			continue
		}

		args = append(args, id)
		query := fmt.Sprintf("UPDATE clientes SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?", strings.Join(updates, ", "))
		if _, err := db.Exec(query, args...); err != nil {
			return stats, err
		}
	}

	return stats, rows.Err()
}

func backfillStudents(db *sql.DB) (backfillStats, error) {
	rows, err := db.Query(`
		SELECT
			id,
			COALESCE(nome, ''),
			COALESCE(nome_protegido, '')
		FROM alunos
	`)
	if err != nil {
		return backfillStats{}, err
	}
	defer rows.Close()

	stats := backfillStats{}

	for rows.Next() {
		var id uint64
		var nome string
		var nomeProtegido string

		if err := rows.Scan(&id, &nome, &nomeProtegido); err != nil {
			return stats, err
		}
		stats.RowsSeen++

		if strings.TrimSpace(nome) == "" || strings.TrimSpace(nomeProtegido) != "" {
			continue
		}

		encrypted, err := security.EncryptPII(strings.TrimSpace(nome))
		if err != nil {
			return stats, fmt.Errorf("alunos.id=%d nome: %w", id, err)
		}

		if _, err := db.Exec(
			"UPDATE alunos SET nome_protegido = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			encrypted,
			id,
		); err != nil {
			return stats, err
		}

		stats.NamesProtected++
	}

	return stats, rows.Err()
}
