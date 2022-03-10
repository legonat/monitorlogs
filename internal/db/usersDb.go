package db

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"time"
)

const (

	WRITE_REFRESH_TOKEN = `INSERT INTO sessions (login, refreshToken, ua, fingerprint, ip, expiresIn, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	GET_REFRESH_SESSION = `SELECT login, refreshToken, fingerprint, expiresIn, createdAt FROM sessions WHERE refreshToken = $1;`

	GET_ALL_REFRESH_SESSIONS = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE login Like $1 ORDER BY createdAt ASC;`

	GET_ALL_REFRESH_SESSIONS_SORTED = `SELECT refreshToken, fingerprint FROM sessions WHERE login Like $1 ORDER BY createdAt ASC;`

	FIND_REFRESH_TOKEN = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE fingerprint = $1;`

	FIND_REFRESH_SESSION = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE refreshToken = $1;`

	DELETE_TOKEN = `DELETE FROM sessions WHERE refreshToken = $1;`

	//INVALIDATE_TOKEN = `UPDATE sessions SET valid = 0 WHERE refreshToken = $1;`

	INSERT_EVENT = `INSERT INTO events (name, caption)VALUES ($1, $2);`

	WRITE_STAT = `INSERT INTO statistics (event, create_at, ip, details) VALUES ($1, $2 ,$3, $4);`

	GET_USER = `SELECT login FROM users WHERE login LIKE $1;`

	WRITE_USER = `INSERT INTO users (login, password, salt, create_at, blocked, try_count, blocked_at, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	GET_USER_PASSWORD = `SELECT password, salt, blocked, try_count FROM users WHERE login LIKE ?;`

	WRITE_TRY_COUNT = `UPDATE users SET try_count = $1 WHERE login LIKE $2 `

	BLOCK_USER = `UPDATE users SET blocked = 1, blocked_at = $1 WHERE login LIKE $2`

	GET_BLOCKED = `SELECT blocked FROM users WHERE login LIKE $1;`

	UNBLOCK_USER = `UPDATE users SET blocked = 0, try_count = 5, blocked_at = 0 WHERE login LIKE $1`

)

type UsersDbSqlite struct {
	db *sql.DB
}

func NewUsersDbSqlite(db *sql.DB) *UsersDbSqlite {
	return &UsersDbSqlite{db: db}
}

func (r *UsersDbSqlite) Register(inputs models.RegisterInputs) error{

	if len(inputs.Password) == 0 {
		return  erx.NewError(604, "Invalid password")
	}

	var user string
	err := r.db.QueryRow(GET_USER, inputs.Login).Scan(&user)
	if err != nil && err!= sql.ErrNoRows {
		return erx.New(err)
	}
	if err == nil {
		_, err = r.db.Exec(WRITE_STAT, 2, GetTime(), inputs.Ip, "User already registered")
		if err != nil {
			return erx.New(err)
		}
		return erx.NewError(603, "User already registered")
	}

	salt, err := GenerateSalt()
	if err != nil {
		return erx.New(err)
	}
	_, err = r.db.Exec(WRITE_USER, inputs.Login, PasswordHash([]byte(inputs.Password), salt), salt, GetTime(), false, 5, 0, false)
	if err != nil {
		return erx.New(err)
	}

	_, err = r.db.Exec(WRITE_STAT, 1, GetTime(), inputs.Ip, "User added")

	return err
}

func (r *UsersDbSqlite) Check(inputs models.LoginInputs) error {

	var user models.User
	err := r.db.QueryRow(GET_USER, inputs.Login).Scan(&user.Login)
	if err != nil && err != sql.ErrNoRows{
		return erx.New(err)
	}

	if err == sql.ErrNoRows{
		_, err = r.db.Exec(WRITE_STAT, 3, GetTime(), inputs.Ip, "User not found")
		return erx.NewError(604, "Invalid password")
	}

	err = r.db.QueryRow(GET_USER_PASSWORD, inputs.Login).Scan(&user.Password, &user.Salt, &user.Blocked, &user.Try_count)
	if err != nil{
		return erx.New(err)
	}

	if user.Blocked {
		fmt.Println("User is blocked")
		_, err = r.db.Exec(WRITE_STAT, 5, GetTime(), inputs.Ip, "Authentication attempt from a blocked user")
		return  erx.NewError(604, "Invalid password")
	}

	passStr := hex.EncodeToString(user.Password)
	checkPass := PasswordHash([]byte(inputs.Password), user.Salt)
	if passStr == hex.EncodeToString(checkPass) {
		fmt.Println("Password is correct")
		_, err = r.db.Exec(WRITE_TRY_COUNT, 5, inputs.Login)
		if err != nil{
			return erx.New(err)
		}
		return nil
	}

	fmt.Println("Password is incorrect")
	user.Try_count--
	_, err = r.db.Exec(WRITE_STAT,4, GetTime(), inputs.Ip, "Invalid password")
	if err != nil{
		return erx.New(err)
	}
	_, err = r.db.Exec(WRITE_TRY_COUNT, user.Try_count, inputs.Login)
	if err != nil{
		return erx.New(err)
	}

	if user.Try_count == 0{
		fmt.Println("User is blocked")
		_, err = r.db.Exec(BLOCK_USER, GetTime(), inputs.Login)
		if err != nil{
			return erx.New(err)
		}

		_, err = r.db.Exec(WRITE_STAT, 6, GetTime(), inputs.Ip, "User blocked")
		if err != nil{
			return erx.New(err)
		}
	}

	return erx.NewError(604,"Invalid password")
}

func (r *UsersDbSqlite) Block(inputs models.BlockInputs) error{

	var blocked bool
	err := r.db.QueryRow(GET_BLOCKED, inputs.Login).Scan(&blocked)
	if err != nil && err != sql.ErrNoRows{
		return erx.New(err)
	}

	if err == sql.ErrNoRows{
		fmt.Println("User not found")
		_, err := r.db.Exec(WRITE_STAT, 3, GetTime(), inputs.Ip, "User not found")
		if err != nil{
			return erx.New(err)
		}
	}

	if blocked == true{
		return erx.NewError(601, "User is already blocked")
	}

	fmt.Println("User is blocked successfully")
	_, err = r.db.Exec(BLOCK_USER, GetTime(), inputs.Login)
	if err != nil{
		return erx.New(err)
	}

	_, err = r.db.Exec(WRITE_STAT, 6, GetTime(), inputs.Ip, "User blocked")
	if err != nil{
		return erx.New(err)
	}

	return err
}

func (r *UsersDbSqlite) Unblock(inputs models.BlockInputs) error{

	var blocked bool
	err := r.db.QueryRow(GET_BLOCKED, inputs.Login).Scan(&blocked)
	if err != nil && err != sql.ErrNoRows{
		return erx.New(err)
	}

	if err == sql.ErrNoRows{
		fmt.Println("User not found")
		_, err := r.db.Exec(WRITE_STAT, 3, GetTime(), inputs.Ip, "User not found")
		if err != nil{
			return erx.New(err)
		}
	}

	if blocked == false {
		return erx.NewError(602, "User is not blocked")
	}

	fmt.Println("User is unblocked successfully")
	_, err = r.db.Exec(UNBLOCK_USER, inputs.Login)
	if err != nil{
		return erx.New(err)
	}

	_, err = r.db.Exec(WRITE_STAT, 7, GetTime(), inputs.Ip, "User unblocked")
	if err != nil{
		return erx.New(err)
	}

	return err
}

func (r *UsersDbSqlite) WriteRefreshToken(inputs models.RefreshSession, daysUntilExpire int) error {

	_, err := r.db.Exec(WRITE_REFRESH_TOKEN, inputs.Login, inputs.Token, inputs.Ua, inputs.Fingerprint, inputs.Ip, GetExpTime(daysUntilExpire), GetTime())
	if err != nil{
		return erx.New(err)
	}

	_, err = r.db.Exec(WRITE_STAT, 8, GetTime(), inputs.Ip, "User auth success")
	if err != nil{
		return erx.New(err)
	}

	return err
}

func (r *UsersDbSqlite) CheckRefreshToken(inputs models.RefreshSession) (string, int, error){
	var login string
	var daysUntilExpire int

	var refSes models.RefreshSession
	err := r.db.QueryRow(GET_REFRESH_SESSION, inputs.Token).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.ExpiresIn, &refSes.CreatedAt)
	if err != nil && err != sql.ErrNoRows{
		return login, daysUntilExpire, erx.New(err)
	}

	if err == sql.ErrNoRows{
		tools.LogWarn("Session not found, Suspicious auth attempt")
		_, err := r.db.Exec(WRITE_STAT, 9, GetTime(), inputs.Ip, "Suspicious auth attempt")
		if err != nil{
			return login, daysUntilExpire, erx.New(err)
		}
		//DeleteAllSessions()
		return login, daysUntilExpire, erx.New(sql.ErrNoRows)
	}

	_, err = r.db.Exec(DELETE_TOKEN, inputs.Token)
	if err != nil{
		return login, daysUntilExpire, erx.New(err)
	}

	if refSes.ExpiresIn < GetTime(){
		tools.LogWarn("Session expired")
		_, err := r.db.Exec(WRITE_STAT, 10, GetTime(), inputs.Ip, "Session expired")
		if err != nil{
			return login, daysUntilExpire, erx.New(err)
		}
		return login, daysUntilExpire, erx.NewError(608, "Session expired")
	}

	if inputs.Fingerprint != refSes.Fingerprint{
		tools.LogWarn("Device not found, Suspicious auth attempt")
		_, err := r.db.Exec(WRITE_STAT, 9, GetTime(), inputs.Ip, "Suspicious auth attempt")
		if err != nil{
			return login, daysUntilExpire, erx.New(err)
		}
		return login, daysUntilExpire, erx.NewError(615, "Suspicious device")
	}

	expiresIn := time.Unix(refSes.ExpiresIn,0)
	createdAt := time.Unix(refSes.CreatedAt, 0)

	daysUntilExpire = int(expiresIn.Sub(createdAt).Hours()) / 24
	login = refSes.Login
	return login,daysUntilExpire, err
}

func (r *UsersDbSqlite) DeleteSession(token string, ip string) error{

	var refSes models.RefreshSession
	err := r.db.QueryRow(FIND_REFRESH_SESSION, token).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.Ip)
	if err != nil && err != sql.ErrNoRows{
		return erx.New(err)
	}

	if err == sql.ErrNoRows{
		return nil
	}

	_, err = r.db.Exec(DELETE_TOKEN, refSes.Token)
	if err != nil{
		return erx.New(err)
	}
	_, err = r.db.Exec(WRITE_STAT, 12, GetTime(), ip, "Session deleted after user request")
	if err != nil{
		return erx.New(err)
	}

	return nil

}

func (r *UsersDbSqlite) TryDeleteOldSession(fingerprint string, ip string) error {

	var refSes models.RefreshSession
	err := r.db.QueryRow(FIND_REFRESH_TOKEN, fingerprint).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.Ip)
	if err != nil && err != sql.ErrNoRows{
		return erx.New(err)
	}

	if err == sql.ErrNoRows{
		return nil
	}

	_, err = r.db.Exec(DELETE_TOKEN, refSes.Token)
	if err != nil{
		return erx.New(err)
	}
	_, err = r.db.Exec(WRITE_STAT, 12, GetTime(), ip, "Session deleted after new login attempt")
	if err != nil{
		return erx.New(err)
	}

	return nil
}



// Function checks count of active Refresh Session. If there are already 3 sessions, func deletes the oldest session
func (r *UsersDbSqlite) CheckSessionsCount(login string, ip string) error {

	rows, err := r.db.Query(GET_ALL_REFRESH_SESSIONS_SORTED, login)
	if err != nil{
		return erx.New(err)
	}
	defer rows.Close()

	var refSession  models.RefreshSession
	var refSlice []models.RefreshSession

	for rows.Next() {
		err = rows.Scan(&refSession.Token, &refSession.Fingerprint)
		if err != nil{
			return erx.New(err)
		}
		refSlice = append(refSlice, refSession)
	}

	if len(refSlice) >= 3 {
		err = r.DeleteSession(refSlice[0].Token, ip)
		if err != nil{
			return erx.New(err)
		}
	}

	return nil

}

func (r *UsersDbSqlite) DeleteAllSessions(inputs models.RefreshSession) error {

	rows, err := r.db.Query(GET_ALL_REFRESH_SESSIONS, inputs.Login)
	if err != nil{
		return erx.New(err)
	}
	defer rows.Close()

	var refSes models.RefreshSession
	var tkns []string
	var sessionValid = false
	for rows.Next() {
		err = rows.Scan(&refSes.Login,&refSes.Token, &refSes.Fingerprint, &refSes.Ip)
		if err != nil{
			return erx.New(err)
		}
		if refSes.Fingerprint == inputs.Fingerprint{
			sessionValid = true
			continue
		}
		tkns = append(tkns, refSes.Token)
	}
	if tkns == nil{
		return erx.NewError(614, "No sessions found")
	}
	if sessionValid {
		for _, v := range tkns {

			_, err := r.db.Exec(DELETE_TOKEN, v)
			if err != nil {
				return erx.New(err)
			}
			_, err = r.db.Exec(WRITE_STAT, 12, GetTime(), inputs.Ip, "Session deleted after Exit Everywhere action")
			if err != nil {
				return erx.New(err)
			}
		}
		return nil
	}

	return erx.NewError(617, "No valid session found")
}