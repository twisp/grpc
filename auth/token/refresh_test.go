package token

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type tokenGeneratorMock struct {
	mock.Mock
}

func (g *tokenGeneratorMock) Generate() ([]byte, error) {
	args := g.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (g *tokenGeneratorMock) TTL() time.Duration {
	args := g.Called()
	return args.Get(0).(time.Duration)
}

func TestRefreshOnStart(t *testing.T) {
	expect := []byte("token")
	now := time.Now()

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(expect, nil)
	generator.On("TTL").Return(time.Minute)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		func() time.Time { return now },
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(time.Second)
	refresher.Stop()

	generated, refreshed, err := refresher.Token()
	require.NoError(t, err)
	require.False(t, refreshed)
	require.Equal(t, expect, generated)

	generator.AssertNumberOfCalls(t, "Generate", 1)
}

func TestRefreshOnStale(t *testing.T) {
	expect := []byte("token")
	now := time.Now()

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(expect, nil)
	generator.On("TTL").Return(time.Minute)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		func() time.Time { return now },
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(time.Second)
	refresher.Stop()

	now = now.Add(time.Hour)

	generated, refreshed, err := refresher.Token()
	require.NoError(t, err)
	require.True(t, refreshed)
	require.Equal(t, expect, generated)

	generator.AssertNumberOfCalls(t, "Generate", 2)
}

func TestRefreshInBackground(t *testing.T) {
	expect := []byte("token")

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(expect, nil)
	generator.On("TTL").Return(time.Second)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		time.Now,
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(1500 * time.Millisecond)
	refresher.Stop()

	for i := 0; i < 10; i++ {
		generated, refreshed, err := refresher.Token()
		require.NoError(t, err)
		require.False(t, refreshed)
		require.Equal(t, expect, generated)
	}

	generator.AssertNumberOfCalls(t, "Generate", 2)
}

func TestRefreshStop(t *testing.T) {
	expect := []byte("token")

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(expect, nil)
	generator.On("TTL").Return(time.Second)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		time.Now,
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(1500 * time.Millisecond)
	refresher.Stop()
	time.Sleep(1500 * time.Millisecond)

	generator.AssertNumberOfCalls(t, "Generate", 2)
}

func TestRefreshAfterBackgroundError(t *testing.T) {
	expect := []byte("token")

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(([]byte)(nil), errors.New("error")).Once()
	generator.On("TTL").Return(time.Second)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		time.Now,
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(time.Second)
	refresher.Stop()

	generator.On("Generate").Return(expect, nil)

	generated, refreshed, err := refresher.Token()
	require.NoError(t, err)
	require.True(t, refreshed)
	require.Equal(t, expect, generated)

	generator.AssertNumberOfCalls(t, "Generate", 2)
}

func TestCachedAfterBackgroundError(t *testing.T) {
	expect := []byte("token")
	now := time.Now()

	generator := new(tokenGeneratorMock)
	generator.On("Generate").Return(expect, nil).Once()
	generator.On("TTL").Return(time.Minute)

	refresher, err := NewTokenRefresherTTL(
		generator,
		1.00,
		1.00,
		func() time.Time { return now },
	)
	require.NoError(t, err)
	require.NotNil(t, refresher)

	time.Sleep(time.Second)
	refresher.Stop()

	generated, refreshed, err := refresher.Token()
	require.NoError(t, err)
	require.False(t, refreshed)
	require.Equal(t, expect, generated)

	now = now.Add(30 * time.Second)

	generated, refreshed, err = refresher.Token()
	require.NoError(t, err)
	require.False(t, refreshed)
	require.Equal(t, expect, generated)

	now = now.Add(30 * time.Second)

	generator.On("Generate").Return(([]byte)(nil), errors.New("error")).Once()

	_, _, err = refresher.Token()
	require.Error(t, err)

	generator.On("Generate").Return(expect, nil)

	generated, refreshed, err = refresher.Token()
	require.NoError(t, err)
	require.True(t, refreshed)
	require.Equal(t, expect, generated)

	generated, refreshed, err = refresher.Token()
	require.NoError(t, err)
	require.False(t, refreshed)
	require.Equal(t, expect, generated)
}
