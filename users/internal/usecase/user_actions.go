package usecase

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"sn/users/internal/domain"
)

func DoesUsernameExist(username string, conn *sql.DB) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);`
	var exists bool
	err := conn.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DoesEmailExist(email string, conn *sql.DB) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`
	var exists bool
	err := conn.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateUser(username, email, password string, conn *sql.DB) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
RETURNING id, created_at, updated_at, last_login;`

	var user domain.User
	err = tx.QueryRow(query, username, email, hashedPassword).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO user_profiles (user_id) VALUES ($1);`
	row := tx.QueryRow(query, user.Id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	return &user, nil
}

func AuthenticateUser(username, password string, conn *sql.DB) (*domain.User, error) {
	query := `SELECT id, password_hash FROM users WHERE username = $1;`
	var user domain.User
	var hashedPassword []byte
	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(query, username).Scan(&user.Id, &hashedPassword)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, err
	}

	query = `UPDATE users SET last_login = NOW() WHERE id = $1;`
	_, err = tx.Exec(query, user.Id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	user.Username = username
	return &user, nil
}

func GetUserProfile(id string, conn *sql.DB) (*domain.User, error) {
	query := `SELECT
    id,
    username,
    email,
    full_name,
    phone_number,
	created_at,
	updated_at,
	last_login,
	is_active
FROM users
WHERE id = $1;`

	var user domain.User
	err := conn.QueryRow(query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsActive)
	if err != nil {
		return nil, err
	}

	query = `SELECT birth_date FROM user_profiles WHERE user_id = $1;`
	err = conn.QueryRow(query, id).Scan(&user.DateOfBirth)

	return &user, nil
}

func UpdateUserProfile(id, fullName, dateOfBirth, phoneNumber string, conn *sql.DB) (*domain.User, error) {
	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if fullName != "" {
		query := `UPDATE users SET full_name = $1 WHERE id = $2;`
		_, err = tx.Exec(query, fullName, id)
		if err != nil {
			return nil, err
		}
	}
	if dateOfBirth != "" {
		query := `UPDATE user_profiles SET birth_date = $1 WHERE user_id = $2;`
		_, err = tx.Exec(query, dateOfBirth, id)
		if err != nil {
			return nil, err
		}
	}
	if phoneNumber != "" {
		query := `UPDATE users SET phone_number = $1 WHERE id = $2;`
		_, err = tx.Exec(query, phoneNumber, id)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	profile, err := GetUserProfile(id, conn)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func DeleteUser(id string, conn *sql.DB) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE id = $1;`
	res, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	query = `DELETE FROM user_profiles WHERE user_id = $1;`
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
