package usersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/tonrock01/another-world-shop/modules/users"
	"github.com/tonrock01/another-world-shop/modules/users/usersPatterns"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
	FindOneOauth(refreshToken string) (*users.Oauth, error)
	UpdateOauth(req *users.UserToken) error
	GetProfile(userId string) (*users.User, error)
	DeleteOauth(oauthId string) error
}

type usersRepository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func UsersRepository(db *sqlx.DB, redisClient *redis.Client) IUsersRepository {
	return &usersRepository{
		db:          db,
		redisClient: redisClient,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	//Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username",
		"role_id"
	FROM "users"
	WHERE "email" = $1;`

	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "oauth"(
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
	RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	fmt.Println("tokem: ", refreshToken)
	query := `
	SELECT
		"id",
		"user_id"
	FROM "oauth"
	WHERE "refresh_token" = $1;`

	oauth := new(users.Oauth)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "id" = :id;`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	ctx := context.Background()

	query := `
	SELECT
		"id",
		"email",
		"username",
		"role_id"
	FROM "users"
	WHERE "id" = $1;`

	profile := new(users.User)

	result, err := r.redisClient.Get(ctx, "user").Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis get failed: %v", err)
	}

	// result, err := r.redisClient.HGetAll(ctx, "user").Result()
	// if err != nil {
	// 	return nil, fmt.Errorf("redis scan failed: %v", err)
	// }
	// if len(result) > 0 {
	// 	// Map ข้อมูลจาก result ไปยัง profile
	// 	profile.Id = result["id"]
	// 	profile.Email = result["email"]
	// 	profile.Username = result["username"]
	// 	profile.RoleId, _ = strconv.Atoi(result["role_id"])
	// 	return profile, nil
	// }

	if result != "" {
		if err := json.Unmarshal([]byte(result), &profile); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %v", err)
		}
		fmt.Printf("Profile from Redis: %+v\n", profile)
		return profile, nil
	} else {
		if err := r.db.Get(profile, query, userId); err != nil {
			return nil, fmt.Errorf("get user failed: %v", err)
		}

		p, err := json.Marshal(&profile)
		if err != nil {
			return nil, fmt.Errorf("marshal failed: %v", err)
		}

		err = r.redisClient.Set(ctx, "user", p, 120*time.Second).Err()
		if err != nil {
			return nil, fmt.Errorf("redis set failed: %v", err)
		}

		// redisHash := map[string]interface{}{
		// 	"id":       profile.Id,
		// 	"email":    profile.Email,
		// 	"username": profile.Username,
		// 	"role_id":  profile.RoleId,
		// }
		// err = r.redisClient.HMSet(ctx, "user", redisHash).Err()
		// if err != nil {
		// 	return nil, fmt.Errorf("redis HMSet failed: %v", err)
		// }
		// r.redisClient.Expire(ctx, "user", 120*time.Second)

		return profile, nil
	}
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	query := `
	DELETE
	FROM "oauth"
	WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil {
		return fmt.Errorf("oauth not found")
	}
	return nil
}
