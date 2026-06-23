package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

type AuditLogsRepository struct {
	db *sql.DB
}

type AuditLogFilters struct {
	Query    string
	User     string
	Action   string
	DateFrom string
	DateTo   string
}

func NewAuditLogsRepository(db *sql.DB) *AuditLogsRepository {
	return &AuditLogsRepository{db: db}
}

func parseAuditJSON(raw []byte) (map[string]interface{}, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func buildAuditChanges(action string, before, after map[string]interface{}) []models.AuditLogFieldChange {
	keysMap := make(map[string]struct{})
	for key := range before {
		keysMap[key] = struct{}{}
	}
	for key := range after {
		keysMap[key] = struct{}{}
	}

	keys := make([]string, 0, len(keysMap))
	for key := range keysMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	changes := make([]models.AuditLogFieldChange, 0, len(keys))
	for _, key := range keys {
		beforeValue, hasBefore := before[key]
		afterValue, hasAfter := after[key]

		switch strings.ToUpper(strings.TrimSpace(action)) {
		case "UPDATE":
			if !hasBefore && !hasAfter {
				continue
			}
			if reflect.DeepEqual(beforeValue, afterValue) {
				continue
			}
		case "INSERT":
			if !hasAfter {
				continue
			}
		case "DELETE":
			if !hasBefore {
				continue
			}
		}

		change := models.AuditLogFieldChange{
			Field: key,
		}
		if hasBefore {
			change.Before = beforeValue
		}
		if hasAfter {
			change.After = afterValue
		}

		changes = append(changes, change)
	}

	return changes
}

func buildAuditSummary(action string, changes []models.AuditLogFieldChange, description string) string {
	normalizedAction := strings.ToUpper(strings.TrimSpace(action))
	if len(changes) == 0 {
		switch normalizedAction {
		case "INSERT":
			return "Registro criado"
		case "DELETE":
			return "Registro removido"
		case "UPDATE":
			if strings.TrimSpace(description) != "" {
				return description
			}
			return "Registro atualizado"
		default:
			if strings.TrimSpace(description) != "" {
				return description
			}
			return normalizedAction
		}
	}

	labels := make([]string, 0, min(3, len(changes)))
	for index, change := range changes {
		if index >= 3 {
			break
		}
		labels = append(labels, change.Field)
	}

	prefix := ""
	switch normalizedAction {
	case "INSERT":
		prefix = "Criado com"
	case "DELETE":
		prefix = "Removido com"
	default:
		prefix = "Alterou"
	}

	if len(changes) <= 3 {
		return fmt.Sprintf("%s %s", prefix, strings.Join(labels, ", "))
	}

	return fmt.Sprintf("%s %d campos: %s", prefix, len(changes), strings.Join(labels, ", "))
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func (r *AuditLogsRepository) buildFiltersQuery(filters AuditLogFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if query := strings.TrimSpace(filters.Query); query != "" {
		likeQuery := "%" + strings.ToLower(query) + "%"
		conditions = append(conditions, `(
			LOWER(COALESCE(u.nome, '')) LIKE ?
			OR LOWER(l.tabela_nome) LIKE ?
			OR LOWER(l.registro_id) LIKE ?
			OR LOWER(COALESCE(l.descricao, '')) LIKE ?
			OR LOWER(CAST(l.dados_antes AS CHAR)) LIKE ?
			OR LOWER(CAST(l.dados_depois AS CHAR)) LIKE ?
		)`)
		args = append(args, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)
	}

	if user := strings.TrimSpace(filters.User); user != "" {
		likeUser := "%" + strings.ToLower(user) + "%"
		conditions = append(conditions, `(LOWER(COALESCE(u.nome, '')) LIKE ? OR CAST(COALESCE(l.id_usuario, 0) AS CHAR) LIKE ?)`)
		args = append(args, likeUser, "%"+user+"%")
	}

	if action := strings.ToUpper(strings.TrimSpace(filters.Action)); action != "" {
		conditions = append(conditions, `UPPER(l.acao) = ?`)
		args = append(args, action)
	}

	if dateFrom := strings.TrimSpace(filters.DateFrom); dateFrom != "" {
		conditions = append(conditions, `DATE(l.created_at) >= ?`)
		args = append(args, dateFrom)
	}

	if dateTo := strings.TrimSpace(filters.DateTo); dateTo != "" {
		conditions = append(conditions, `DATE(l.created_at) <= ?`)
		args = append(args, dateTo)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *AuditLogsRepository) FetchPage(page, pageSize int, filters AuditLogFilters) (models.AuditLogPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	whereClause, args := r.buildFiltersQuery(filters)

	var total uint64
	countQuery := `
		SELECT COUNT(*)
		FROM treehousedb.logs_auditoria l
		LEFT JOIN treehousedb.usuarios u
			ON u.id = l.id_usuario
	` + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return models.AuditLogPage{}, err
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT
			l.id,
			l.id_usuario,
			COALESCE(u.nome, ''),
			l.tabela_nome,
			l.registro_id,
			l.acao,
			CAST(l.dados_antes AS CHAR),
			CAST(l.dados_depois AS CHAR),
			COALESCE(l.descricao, ''),
			COALESCE(l.ip_origem, ''),
			COALESCE(l.user_agent, ''),
			l.created_at
		FROM treehousedb.logs_auditoria l
		LEFT JOIN treehousedb.usuarios u
			ON u.id = l.id_usuario
	` + whereClause + `
		ORDER BY l.created_at DESC, l.id DESC
		LIMIT ? OFFSET ?
	`

	queryArgs := append(append([]interface{}{}, args...), pageSize, offset)
	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return models.AuditLogPage{}, err
	}
	defer rows.Close()

	items := make([]models.AuditLogItem, 0, pageSize)
	for rows.Next() {
		var item models.AuditLogItem
		var userID sql.NullInt64
		var beforeRaw sql.NullString
		var afterRaw sql.NullString

		if err = rows.Scan(
			&item.ID,
			&userID,
			&item.UserName,
			&item.TableName,
			&item.RecordID,
			&item.Action,
			&beforeRaw,
			&afterRaw,
			&item.Description,
			&item.IPOrigin,
			&item.UserAgent,
			&item.CreatedAt,
		); err != nil {
			return models.AuditLogPage{}, err
		}

		if userID.Valid {
			value := uint64(userID.Int64)
			item.UserID = &value
		}
		if strings.TrimSpace(item.UserName) == "" {
			item.UserName = "Sistema"
		}

		if beforeRaw.Valid {
			item.Before, err = parseAuditJSON([]byte(beforeRaw.String))
			if err != nil {
				return models.AuditLogPage{}, err
			}
		}
		if afterRaw.Valid {
			item.After, err = parseAuditJSON([]byte(afterRaw.String))
			if err != nil {
				return models.AuditLogPage{}, err
			}
		}

		item.Changes = buildAuditChanges(item.Action, item.Before, item.After)
		item.ChangesCount = len(item.Changes)
		item.Summary = buildAuditSummary(item.Action, item.Changes, item.Description)

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return models.AuditLogPage{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + uint64(pageSize) - 1) / uint64(pageSize))
	}

	return models.AuditLogPage{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
