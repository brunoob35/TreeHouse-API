package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

// UsersRepository is responsible for all database operations related to users.
type UsersRepository struct {
	db          *sql.DB
	auditUserID *uint64
}

type userRowScanner interface {
	Scan(dest ...interface{}) error
}

// NewUsersRepository creates a new repository instance for users.
func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) WithAuditUser(userID uint64) *UsersRepository {
	r.auditUserID = &userID
	return r
}

func nullableUserString(value string) interface{} {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func userPasswordValue(value string) string {
	return strings.TrimSpace(value)
}

func encryptUserPII(value string) (interface{}, interface{}, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil, nil
	}

	encrypted, err := security.EncryptPII(trimmed)
	if err != nil {
		return nil, nil, err
	}

	hash, err := security.HashPII(trimmed)
	if err != nil {
		return nil, nil, err
	}

	return encrypted, hash, nil
}

func scanUser(scanner userRowScanner) (models.User, error) {
	var user models.User
	var password sql.NullString
	var cpfLegacy sql.NullString
	var rgLegacy sql.NullString
	var cpfProtected sql.NullString
	var rgProtected sql.NullString
	var phone sql.NullString
	var birth sql.NullTime
	var lgpdAcceptedAt sql.NullTime
	var lgpdPurpose sql.NullString

	err := scanner.Scan(
		&user.ID,
		&user.IDEndereco,
		&password,
		&user.Nome,
		&user.Email,
		&cpfLegacy,
		&rgLegacy,
		&cpfProtected,
		&rgProtected,
		&phone,
		&user.Ativo,
		&birth,
		&user.LGPDAceito,
		&lgpdAcceptedAt,
		&lgpdPurpose,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	if password.Valid {
		user.Senha = password.String
	}
	if cpfProtected.Valid {
		user.CPF, err = security.DecryptPII(cpfProtected.String)
		if err != nil {
			return models.User{}, err
		}
	} else if cpfLegacy.Valid {
		user.CPF = cpfLegacy.String
	}
	if rgProtected.Valid {
		user.RG, err = security.DecryptPII(rgProtected.String)
		if err != nil {
			return models.User{}, err
		}
	} else if rgLegacy.Valid {
		user.RG = rgLegacy.String
	}
	if phone.Valid {
		user.Telefone = phone.String
	}
	if birth.Valid {
		user.Nascimento = &birth.Time
	}
	if lgpdAcceptedAt.Valid {
		user.LGPDAceitoEm = &lgpdAcceptedAt.Time
	}
	if lgpdPurpose.Valid {
		user.LGPDFinalidade = lgpdPurpose.String
	}

	return user, nil
}

// FetchByID searches for a user by its ID.
// It returns the user base data without loading permission relations.
// Permissions are intentionally loaded separately because the application now
// aggregates them into a numeric bitmask for authentication and authorization.
func (r *UsersRepository) FetchByID(id uint64) (models.User, error) {
	query := `
		SELECT
			id,
			id_endereco,
			senha,
			nome,
			email,
			cpf,
			rg,
			cpf_protegido,
			rg_protegido,
			telefone,
			ativo,
			nascimento,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade,
			created_at,
			updated_at
		FROM treehousedb.usuarios
		WHERE id = ?
	`

	user, err := scanUser(r.db.QueryRow(query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, sql.ErrNoRows
		}
		return models.User{}, err
	}

	return user, nil
}

// FetchByEmail searches for a user by email.
// This function is usually used during the login flow before validating
// the password and generating the JWT token.
func (r *UsersRepository) FetchByEmail(email string) (models.User, error) {
	query := `
		SELECT
			id,
			id_endereco,
			senha,
			nome,
			email,
			cpf,
			rg,
			cpf_protegido,
			rg_protegido,
			telefone,
			ativo,
			nascimento,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade,
			created_at,
			updated_at
		FROM treehousedb.usuarios
		WHERE email = ?
	`

	user, err := scanUser(r.db.QueryRow(query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, sql.ErrNoRows
		}
		return models.User{}, err
	}

	return user, nil
}

// FetchPermissionIDsByUser returns all permission IDs assigned to a user.
// The permission IDs stored in the database must already be valid bit flags
// (1, 2, 4, 8, 16, ...). Because of that, they can later be aggregated into
// a single numeric mask using the bitwise OR operator.
func (r *UsersRepository) FetchPermissionIDsByUser(userID uint64) ([]uint64, error) {
	query := `
		SELECT up.id_permissao
		FROM treehousedb.usuarios_permissoes up
		WHERE up.id_usuario = ?
		ORDER BY up.id_permissao
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []uint64

	for rows.Next() {
		var permissionID uint64

		if err = rows.Scan(&permissionID); err != nil {
			return nil, err
		}

		permissions = append(permissions, permissionID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// FetchPermissionMaskByUser loads all permission IDs of a user and
// aggregates them into a single numeric bitmask.
// Example:
// If the user has permissions [1, 4], the final mask will be:
//
//	1 | 4 = 5
//
// This mask is the value that should be stored in the JWT token.
func (r *UsersRepository) FetchPermissionMaskByUser(userID uint64) (uint64, error) {
	permissionIDs, err := r.FetchPermissionIDsByUser(userID)
	if err != nil {
		return 0, err
	}

	return authentication.BuildPermissionMask(permissionIDs), nil
}

// Insert creates a new user record in the database.
// This function inserts only the user base record. Permission assignments
// must be handled separately through the relation table "usuarios_permissoes".
func (r *UsersRepository) Insert(user models.User) (uint64, error) {
	cpfProtected, cpfHash, err := encryptUserPII(user.CPF)
	if err != nil {
		return 0, err
	}

	rgProtected, rgHash, err := encryptUserPII(user.RG)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return 0, err
	}

	query := `
		INSERT INTO treehousedb.usuarios (
			senha,
			nome,
			email,
			cpf,
			rg,
			cpf_protegido,
			rg_protegido,
			cpf_hash,
			rg_hash,
			telefone,
			ativo,
			nascimento,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		userPasswordValue(user.Senha),
		user.Nome,
		user.Email,
		nil,
		nil,
		cpfProtected,
		rgProtected,
		cpfHash,
		rgHash,
		nullableUserString(user.Telefone),
		user.Ativo,
		user.Nascimento,
		user.LGPDAceito,
		user.LGPDAceitoEm,
		nullableUserString(user.LGPDFinalidade),
	)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(insertedID), nil
}

// InsertWithPermission creates a new user and associates a permission
// in the same transaction.
func (r *UsersRepository) InsertWithPermission(user models.User, permissionID uint64) (uint64, error) {
	cpfProtected, cpfHash, err := encryptUserPII(user.CPF)
	if err != nil {
		return 0, err
	}

	rgProtected, rgHash, err := encryptUserPII(user.RG)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	insertUserQuery := `
		INSERT INTO treehousedb.usuarios (
			senha,
			nome,
			email,
			cpf,
			rg,
			cpf_protegido,
			rg_protegido,
			cpf_hash,
			rg_hash,
			telefone,
			ativo,
			nascimento,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		insertUserQuery,
		userPasswordValue(user.Senha),
		user.Nome,
		user.Email,
		nil,
		nil,
		cpfProtected,
		rgProtected,
		cpfHash,
		rgHash,
		nullableUserString(user.Telefone),
		user.Ativo,
		user.Nascimento,
		user.LGPDAceito,
		user.LGPDAceitoEm,
		nullableUserString(user.LGPDFinalidade),
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	insertPermissionQuery := `
		INSERT INTO treehousedb.usuarios_permissoes (
			id_usuario,
			id_permissao
		) VALUES (?, ?)
	`

	_, err = tx.Exec(insertPermissionQuery, insertedID, permissionID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("usuario criado, mas falha ao associar permissao: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(insertedID), nil
}

// Update updates the user base data.
// Permission assignments are not updated here because they are stored in a
// separate many-to-many relation table.
func (r *UsersRepository) Update(id uint64, user models.User) error {
	cpfProtected, cpfHash, err := encryptUserPII(user.CPF)
	if err != nil {
		return err
	}

	rgProtected, rgHash, err := encryptUserPII(user.RG)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return err
	}

	query := `
		UPDATE treehousedb.usuarios
		SET
			id_endereco = ?,
			nome = ?,
			email = ?,
			cpf = ?,
			rg = ?,
			cpf_protegido = ?,
			rg_protegido = ?,
			cpf_hash = ?,
			rg_hash = ?,
			telefone = ?,
			ativo = ?,
			nascimento = ?,
			lgpd_aceito = ?,
			lgpd_aceito_em = ?,
			lgpd_finalidade = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(
		query,
		user.IDEndereco,
		user.Nome,
		user.Email,
		nil,
		nil,
		cpfProtected,
		rgProtected,
		cpfHash,
		rgHash,
		nullableUserString(user.Telefone),
		user.Ativo,
		user.Nascimento,
		user.LGPDAceito,
		user.LGPDAceitoEm,
		nullableUserString(user.LGPDFinalidade),
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum user encontrado com id %d", id)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates only the user's password hash.
func (r *UsersRepository) UpdatePassword(userID uint64, senhaHash string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return err
	}

	statement, err := tx.Prepare(`
		UPDATE treehousedb.usuarios
		SET senha = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(senhaHash, userID); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete performs a soft delete on a user.
// Instead of removing the row from the database,
// this operation sets the "ativo" field to false.
func (r *UsersRepository) Delete(id uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return err
	}

	query := `
		UPDATE treehousedb.usuarios
		SET
			ativo = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum user encontrado com id %d", id)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// FetchAllUsers returns a list of users optionally filtered by name.
// The "nome" parameter is optional. If provided, the query will perform
// a case-insensitive search using a LIKE clause.
func (r *UsersRepository) FetchAllUsers(nome string) ([]models.User, error) {
	query := `
		SELECT
			id,
			id_endereco,
			senha,
			nome,
			email,
			cpf,
			rg,
			cpf_protegido,
			rg_protegido,
			telefone,
			ativo,
			nascimento,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade,
			created_at,
			updated_at
		FROM treehousedb.usuarios
	`

	var args []interface{}

	if nome != "" {
		query += " WHERE LOWER(nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += " ORDER BY nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		user, scanErr := scanUser(rows)
		err = scanErr
		if err != nil {
			return nil, err
		}

		// Never expose password hashes in API responses
		user.Senha = ""

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FetchAllActiveUsers returns all active users optionally filtered by name.
func (r *UsersRepository) FetchAllActiveUsers(nome string) ([]models.User, error) {
	query := `
		SELECT
			id,
			id_endereco,
			senha,
			nome,
			email,
			cpf,
			rg,
			telefone,
			ativo,
			nascimento,
			created_at,
			updated_at
		FROM treehousedb.usuarios
		WHERE ativo = TRUE
	`

	var args []interface{}

	if nome != "" {
		query += " AND LOWER(nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += " ORDER BY nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		user, scanErr := scanUser(rows)
		err = scanErr
		if err != nil {
			return nil, err
		}

		user.Senha = ""
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FetchProfessors returns all active users with professor permission optionally filtered by name.
func (r *UsersRepository) FetchProfessors(nome string) ([]models.User, error) {
	log.Println("Beteu no Repo")
	query := `
		SELECT
			u.id,
			u.id_endereco,
			u.senha,
			u.nome,
			u.email,
			u.cpf,
			u.rg,
			u.cpf_protegido,
			u.rg_protegido,
			u.telefone,
			u.ativo,
			u.nascimento,
			u.lgpd_aceito,
			u.lgpd_aceito_em,
			u.lgpd_finalidade,
			u.created_at,
			u.updated_at
		FROM treehousedb.usuarios u
		INNER JOIN treehousedb.usuarios_permissoes up
			ON up.id_usuario = u.id
		WHERE up.id_permissao = 2
		  AND u.ativo = TRUE
	`

	var args []interface{}

	if nome != "" {
		query += " AND LOWER(u.nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += " ORDER BY u.nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		user, scanErr := scanUser(rows)
		err = scanErr
		if err != nil {
			return nil, err
		}

		user.Senha = ""
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ReturnAllProfessors returns all users with professor permission optionally filtered by name.
func (r *UsersRepository) ReturnAllProfessors(nome string) ([]models.User, error) {
	query := `
		SELECT
			u.id,
			u.id_endereco,
			u.senha,
			u.nome,
			u.email,
			u.cpf,
			u.rg,
			u.cpf_protegido,
			u.rg_protegido,
			u.telefone,
			u.ativo,
			u.nascimento,
			u.lgpd_aceito,
			u.lgpd_aceito_em,
			u.lgpd_finalidade,
			u.created_at,
			u.updated_at
		FROM treehousedb.usuarios u
		INNER JOIN treehousedb.usuarios_permissoes up
			ON up.id_usuario = u.id
		WHERE up.id_permissao = 2
	`

	var args []interface{}

	if nome != "" {
		query += " AND LOWER(u.nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += " ORDER BY u.nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		user, scanErr := scanUser(rows)
		err = scanErr
		if err != nil {
			return nil, err
		}

		user.Senha = ""
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) CountClassesByProfessorIDs(professorIDs []uint64) ([]models.ProfessorClassCount, error) {
	if len(professorIDs) == 0 {
		return []models.ProfessorClassCount{}, nil
	}

	placeholders := make([]string, len(professorIDs))
	args := make([]interface{}, 0, len(professorIDs)+1)

	for i, id := range professorIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			u.id AS professor_id,
			COUNT(t.id) AS classes_count
		FROM treehousedb.usuarios u
		INNER JOIN treehousedb.usuarios_permissoes up
			ON up.id_usuario = u.id
			AND up.id_permissao = ?
		LEFT JOIN treehousedb.turmas t
			ON t.id_professor = u.id
			AND t.deleted_at IS NULL
		WHERE u.id IN (%s)
		GROUP BY u.id
		ORDER BY u.id
	`, strings.Join(placeholders, ","))

	args = append([]interface{}{uint64(authentication.PermProfessor)}, args...)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []models.ProfessorClassCount

	for rows.Next() {
		var item models.ProfessorClassCount

		if err = rows.Scan(&item.ProfessorID, &item.ClassesCount); err != nil {
			return nil, err
		}

		counts = append(counts, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}
