package user

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type UserProfileRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    UserProfileRepository
	cleanup func()
}

func (s *UserProfileRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "test_db",
				"POSTGRES_USER":     "test_user",
				"POSTGRES_PASSWORD": "test_pass",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10 * time.Second),
		},
		Started: true,
	})
	require.NoError(s.T(), err, "启动 PostgreSQL 容器失败")

	host, err := container.Host(ctx)
	require.NoError(s.T(), err, "获取容器主机地址失败")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(s.T(), err, "获取容器映射端口失败")

	dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_pass dbname=test_db sslmode=disable",
		host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err, "连接数据库失败")

	sqlDB, err := db.DB()
	require.NoError(s.T(), err, "获取 SQL DB 失败")

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = db.AutoMigrate(&models.UserProfile{})
	require.NoError(s.T(), err, "数据库迁移失败")

	s.db = db
	s.repo = NewUserProfileRepository(db)
	s.cleanup = func() {
		_ = container.Terminate(ctx)
	}
}

func (s *UserProfileRepositoryTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func (s *UserProfileRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM user_profiles")
}

func (s *UserProfileRepositoryTestSuite) TearDownTest() {
}

func (s *UserProfileRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		RealName:     "测试用户",
		Status:       models.UserStatusActive,
	}

	err := s.repo.Create(ctx, user)

	require.NoError(s.T(), err)
	assert.NotZero(s.T(), user.ID)
	assert.Equal(s.T(), "testuser", user.Username)
	assert.Equal(s.T(), "test@example.com", user.Email)
}

func (s *UserProfileRepositoryTestSuite) TestCreate_Validation() {
	tests := []struct {
		name    string
		user    *models.UserProfile
		wantErr bool
		errMsg  string
	}{
		{
			name: "用户名为空",
			user: &models.UserProfile{
				PasswordHash: "hash",
				Email:        "test@example.com",
			},
			wantErr: true,
			errMsg:  "用户名不能为空",
		},
		{
			name: "用户名太短",
			user: &models.UserProfile{
				Username:     "ab",
				PasswordHash: "hash",
			},
			wantErr: true,
			errMsg:  "用户名长度必须在3-20个字符之间",
		},
		{
			name: "用户名太长",
			user: &models.UserProfile{
				Username:     "this_username_is_way_too_long_for_the_limit",
				PasswordHash: "hash",
			},
			wantErr: true,
			errMsg:  "用户名长度必须在3-20个字符之间",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.repo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *UserProfileRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByID(ctx, user.ID.String())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), user.ID, found.ID)
	assert.Equal(s.T(), "testuser", found.Username)
}

func (s *UserProfileRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestGetByUsername_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByUsername(ctx, "testuser")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "testuser", found.Username)
}

func (s *UserProfileRepositoryTestSuite) TestGetByUsername_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByUsername(ctx, "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestGetByEmail_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByEmail(ctx, "test@example.com")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "test@example.com", found.Email)
}

func (s *UserProfileRepositoryTestSuite) TestGetByEmail_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByEmail(ctx, "nonexistent@example.com")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestGetByEmail_Empty() {
	ctx := context.Background()

	found, err := s.repo.GetByEmail(ctx, "")

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestExistsByID_True() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.ExistsByID(ctx, user.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestExistsByID_False() {
	ctx := context.Background()

	exists, err := s.repo.ExistsByID(ctx, "00000000-0000-0000-0000-000000000000")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckUsernameExists_True() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckUsernameExists(ctx, "testuser")

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckUsernameExists_False() {
	ctx := context.Background()

	exists, err := s.repo.CheckUsernameExists(ctx, "nonexistent")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckUsernameExists_ExcludeID() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckUsernameExists(ctx, "testuser", user.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckEmailExists_True() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckEmailExists(ctx, "test@example.com")

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckEmailExists_False() {
	ctx := context.Background()

	exists, err := s.repo.CheckEmailExists(ctx, "nonexistent@example.com")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckEmailExists_Empty() {
	ctx := context.Background()

	exists, err := s.repo.CheckEmailExists(ctx, "")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestSetMustChangePassword_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:           "testuser",
		PasswordHash:       "hashedpassword",
		MustChangePassword: false,
		Status:             models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.SetMustChangePassword(ctx, user.ID.String(), true)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.True(s.T(), updated.MustChangePassword)
}

func (s *UserProfileRepositoryTestSuite) TestUpdatePassword_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:           "testuser",
		PasswordHash:       "oldhash",
		MustChangePassword: true,
		LoginAttempts:      5,
		Status:             models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.UpdatePassword(ctx, user.ID.String(), "newhash")

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "newhash", updated.PasswordHash)
	assert.False(s.T(), updated.MustChangePassword)
	assert.Equal(s.T(), int32(0), updated.LoginAttempts)
}

func (s *UserProfileRepositoryTestSuite) TestUpdatePassword_UserNotFound() {
	ctx := context.Background()

	err := s.repo.UpdatePassword(ctx, "00000000-0000-0000-0000-000000000000", "newhash")

	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "old@example.com",
		RealName:     "旧名称",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	user.Email = "new@example.com"
	user.RealName = "新名称"

	err = s.repo.Update(ctx, user)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "new@example.com", updated.Email)
	assert.Equal(s.T(), "新名称", updated.RealName)
}

func (s *UserProfileRepositoryTestSuite) TestUpdate_SystemUser_ChangeUsername() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "systemuser",
		PasswordHash: "hashedpassword",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	user.Username = "newusername"

	err = s.repo.Update(ctx, user)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "系统用户的用户名不能被修改")
}

func (s *UserProfileRepositoryTestSuite) TestSoftDelete_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.SoftDelete(ctx, user.ID.String())

	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(ctx, user.ID.String())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
}

func (s *UserProfileRepositoryTestSuite) TestSoftDelete_SystemUser() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "systemuser",
		PasswordHash: "hashedpassword",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.SoftDelete(ctx, user.ID.String())

	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestIncrementLoginAttempts_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:      "testuser",
		PasswordHash:  "hashedpassword",
		LoginAttempts: 2,
		Status:        models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.IncrementLoginAttempts(ctx, user.ID.String())

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int32(3), updated.LoginAttempts)
}

func (s *UserProfileRepositoryTestSuite) TestResetLoginAttempts_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:      "testuser",
		PasswordHash:  "hashedpassword",
		LoginAttempts: 5,
		Status:        models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.ResetLoginAttempts(ctx, user.ID.String())

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int32(0), updated.LoginAttempts)
}

func (s *UserProfileRepositoryTestSuite) TestUpdateLastLoginTime_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.UpdateLastLoginTime(ctx, user.ID.String())
	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), updated.LastLoginTime)

	assert.True(s.T(), *updated.LastLoginTime > 0)
}

func (s *UserProfileRepositoryTestSuite) TestFindSystemUsers() {
	ctx := context.Background()

	systemUser1 := &models.UserProfile{
		Username:     "system1",
		PasswordHash: "hash1",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	systemUser2 := &models.UserProfile{
		Username:     "system2",
		PasswordHash: "hash2",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	normalUser := &models.UserProfile{
		Username:     "normal",
		PasswordHash: "hash3",
		IsSystemUser: false,
		Status:       models.UserStatusActive,
	}

	err := s.repo.Create(ctx, systemUser1)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, systemUser2)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, normalUser)
	require.NoError(s.T(), err)

	users, err := s.repo.FindSystemUsers(ctx)

	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)

	usernames := make([]string, len(users))
	for i, u := range users {
		usernames[i] = u.Username
	}

	assert.Contains(s.T(), usernames, "system1")
	assert.Contains(s.T(), usernames, "system2")
	assert.NotContains(s.T(), usernames, "normal")
}

func (s *UserProfileRepositoryTestSuite) TestIsSystemUser_True() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "systemuser",
		PasswordHash: "hash",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	isSystem, err := s.repo.IsSystemUser(ctx, user.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), isSystem)
}

func (s *UserProfileRepositoryTestSuite) TestIsSystemUser_False() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "normaluser",
		PasswordHash: "hash",
		IsSystemUser: false,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	isSystem, err := s.repo.IsSystemUser(ctx, user.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), isSystem)
}

func (s *UserProfileRepositoryTestSuite) TestGetByPhone_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Phone:        "13800138000",
		Email:        "test@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByPhone(ctx, "13800138000")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "13800138000", found.Phone)
}

func (s *UserProfileRepositoryTestSuite) TestGetByPhone_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByPhone(ctx, "19900000000")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestGetByPhone_Empty() {
	ctx := context.Background()

	found, err := s.repo.GetByPhone(ctx, "")

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestUpdateLoginAttempts_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:      "testuser",
		PasswordHash:  "hashedpassword",
		LoginAttempts: 2,
		Status:        models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.UpdateLoginAttempts(ctx, user.ID.String(), 5)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int32(5), updated.LoginAttempts)
}

func (s *UserProfileRepositoryTestSuite) TestUpdateLoginAttempts_UserNotFound() {
	ctx := context.Background()

	err := s.repo.UpdateLoginAttempts(ctx, "00000000-0000-0000-0000-000000000000", 3)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "用户不存在或已删除")
}

func (s *UserProfileRepositoryTestSuite) TestCheckPhoneExists_True() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Phone:        "13800138000",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckPhoneExists(ctx, "13800138000")

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckPhoneExists_False() {
	ctx := context.Background()

	exists, err := s.repo.CheckPhoneExists(ctx, "19900000000")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckPhoneExists_Empty() {
	ctx := context.Background()

	exists, err := s.repo.CheckPhoneExists(ctx, "")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestCheckPhoneExists_ExcludeID() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Phone:        "13800138000",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckPhoneExists(ctx, "13800138000", user.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserProfileRepositoryTestSuite) TestFindByMedicalLicense_Success() {
	ctx := context.Background()

	license := "LICENSE123"
	user := &models.UserProfile{
		Username:      "testuser",
		PasswordHash:  "hashedpassword",
		LicenseNumber: license,
		Status:        models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	found, err := s.repo.FindByMedicalLicense(ctx, license)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), license, found.LicenseNumber)
}

func (s *UserProfileRepositoryTestSuite) TestFindByMedicalLicense_NotFound() {
	ctx := context.Background()

	found, err := s.repo.FindByMedicalLicense(ctx, "NONEXISTENT")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestFindByMedicalLicense_Empty() {
	ctx := context.Background()

	found, err := s.repo.FindByMedicalLicense(ctx, "")

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestFindBySpecialty_Success() {
	ctx := context.Background()

	user1 := &models.UserProfile{
		Username:     "user1",
		PasswordHash: "hash1",
		Specialties:  `["Cardiology", "Internal Medicine"]`,
		Status:       models.UserStatusActive,
	}
	user2 := &models.UserProfile{
		Username:     "user2",
		PasswordHash: "hash2",
		Specialties:  `["Cardiology", "Surgery"]`,
		Status:       models.UserStatusActive,
	}
	user3 := &models.UserProfile{
		Username:     "user3",
		PasswordHash: "hash3",
		Specialties:  `["Pediatrics"]`,
		Status:       models.UserStatusActive,
	}

	err := s.repo.Create(ctx, user1)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, user2)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, user3)
	require.NoError(s.T(), err)

	users, _, err := s.repo.FindBySpecialty(ctx, "Cardiology", nil)

	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)

	usernames := make([]string, len(users))
	for i, u := range users {
		usernames[i] = u.Username
	}

	assert.Contains(s.T(), usernames, "user1")
	assert.Contains(s.T(), usernames, "user2")
	assert.NotContains(s.T(), usernames, "user3")
}

func (s *UserProfileRepositoryTestSuite) TestFindBySpecialty_NotFound() {
	ctx := context.Background()

	users, _, err := s.repo.FindBySpecialty(ctx, "Nonexistent", nil)

	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 0)
}

func (s *UserProfileRepositoryTestSuite) TestHardDelete_Success() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		IsSystemUser: false,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.HardDelete(ctx, user.ID.String())

	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(ctx, user.ID.String())
	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestHardDelete_SystemUser() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "systemuser",
		PasswordHash: "hashedpassword",
		IsSystemUser: true,
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	err = s.repo.HardDelete(ctx, user.ID.String())

	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestHardDelete_UserNotFound() {
	ctx := context.Background()

	err := s.repo.HardDelete(ctx, "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestFindWithConditions_StatusFilter() {
	ctx := context.Background()

	// 创建不同状态的用户
	user1 := &models.UserProfile{
		Username:     "active_user",
		PasswordHash: "hash1",
		Status:       models.UserStatusActive,
	}
	user2 := &models.UserProfile{
		Username:     "locked_user",
		PasswordHash: "hash2",
		Status:       models.UserStatusLocked,
	}

	err := s.repo.Create(ctx, user1)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, user2)
	require.NoError(s.T(), err)

	// 测试只查询活跃用户
	status := models.UserStatusActive
	conditions := &UserProfileQueryConditions{
		Status: &status,
		Page: &base.QueryOptions{
			Page:     1,
			PageSize: 10,
		},
	}
	users, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(users), 1)

	hasActiveUser := false

	for _, u := range users {
		if u.Username == "active_user" {
			hasActiveUser = true
			break
		}
	}

	assert.True(s.T(), hasActiveUser, "应该包含活跃用户")
}

func (s *UserProfileRepositoryTestSuite) TestFindWithConditions_EmailSearch() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "searchuser",
		PasswordHash: "hash",
		Email:        "search@example.com",
		Status:       models.UserStatusActive,
	}
	err := s.repo.Create(ctx, user)
	require.NoError(s.T(), err)

	conditions := &UserProfileQueryConditions{
		Page: &base.QueryOptions{
			Search:   "search",
			Page:     1,
			PageSize: 10,
		},
	}
	users, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.Greater(s.T(), len(users), 0)

	found := false

	for _, u := range users {
		if u.Username == "searchuser" {
			found = true
			break
		}
	}

	assert.True(s.T(), found)
}

func (s *UserProfileRepositoryTestSuite) TestFindWithConditions_Pagination() {
	ctx := context.Background()

	// 创建多个用户
	for i := 1; i <= 5; i++ {
		user := &models.UserProfile{
			Username:     fmt.Sprintf("pageuser%d", i),
			PasswordHash: "hash",
			Status:       models.UserStatusActive,
		}
		err := s.repo.Create(ctx, user)
		require.NoError(s.T(), err)
	}

	// 测试分页
	conditions := &UserProfileQueryConditions{
		Page: &base.QueryOptions{
			Page:     1,
			PageSize: 2,
		},
	}
	users, total, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(users), 2)
	assert.GreaterOrEqual(s.T(), total, int64(5))
}

func (s *UserProfileRepositoryTestSuite) TestFindWithConditions_Sorting() {
	ctx := context.Background()

	user1 := &models.UserProfile{
		Username:     "zebra",
		PasswordHash: "hash1",
		Status:       models.UserStatusActive,
	}
	user2 := &models.UserProfile{
		Username:     "apple",
		PasswordHash: "hash2",
		Status:       models.UserStatusActive,
	}

	err := s.repo.Create(ctx, user1)
	require.NoError(s.T(), err)
	err = s.repo.Create(ctx, user2)
	require.NoError(s.T(), err)

	// 测试按用户名升序排序
	conditions := &UserProfileQueryConditions{
		Page: &base.QueryOptions{
			Page:      1,
			PageSize:  10,
			OrderBy:   "username",
			OrderDesc: false,
		},
	}
	users, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(users), 2)

	// 检查是否按升序排列
	for i := 1; i < len(users); i++ {
		assert.LessOrEqual(s.T(), users[i-1].Username, users[i].Username)
	}
}

func (s *UserProfileRepositoryTestSuite) TestWithTx_TransactionRollback() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "txuser",
		PasswordHash: "hash",
		Status:       models.UserStatusActive,
	}

	// 开始事务
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoWithTx := s.repo.WithTx(tx)

		// 在事务中创建用户
		err := repoWithTx.Create(ctx, user)
		require.NoError(s.T(), err)

		// 用户应该存在于事务中
		found, err := repoWithTx.GetByID(ctx, user.ID.String())
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), found)

		// 回滚事务
		return errors.New("rollback transaction")
	})

	assert.Error(s.T(), err)

	// 用户不应该存在于数据库中
	_, err = s.repo.GetByID(ctx, user.ID.String())
	assert.Error(s.T(), err)
}

func (s *UserProfileRepositoryTestSuite) TestWithTx_TransactionCommit() {
	ctx := context.Background()

	user := &models.UserProfile{
		Username:     "txcommituser",
		PasswordHash: "hash",
		Status:       models.UserStatusActive,
	}

	// 开始事务
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoWithTx := s.repo.WithTx(tx)

		// 在事务中创建用户
		err := repoWithTx.Create(ctx, user)
		if err != nil {
			return err
		}

		// 用户应该存在于事务中
		found, err := repoWithTx.GetByID(ctx, user.ID.String())
		if err != nil {
			return err
		}

		if found == nil {
			return errors.New("用户未找到")
		}

		// 提交事务（返回 nil）
		return nil
	})

	require.NoError(s.T(), err)

	// 用户应该存在于数据库中
	found, err := s.repo.GetByID(ctx, user.ID.String())
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), "txcommituser", found.Username)
}

func TestUserProfileRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserProfileRepositoryTestSuite))
}
