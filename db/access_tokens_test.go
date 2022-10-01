//go:build integration

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessTokensQueries(t *testing.T) {
	ctx := context.Background()

	withTestDB(t, func(db DBTX) {
		querier := New()

		var (
			validToken         IssueAccessTokenRow
			validNoExpiryToken IssueAccessTokenRow
			expiredToken       IssueAccessTokenRow
			revokedToken       IssueAccessTokenRow
		)
		t.Run("IssueAccessToken", func(t *testing.T) {
			var err error
			validToken, err = querier.IssueAccessToken(ctx, db, IssueAccessTokenParams{
				Name:               "valid",
				TokenHash:          "valid-hash",
				NamespaceSelectors: []string{"ns-0", "ns-1"},
				Scopes:             []string{string(AccessTokenScopeAdmin), string(AccessTokenScopeControlRw)},
				ExpiresAt:          sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
			})
			require.NoError(t, err)

			validNoExpiryToken, err = querier.IssueAccessToken(ctx, db, IssueAccessTokenParams{
				Name:               "valid no expiry",
				TokenHash:          "valid-no-expriy-hash",
				NamespaceSelectors: []string{"ns-0", "ns-1"},
				Scopes:             []string{string(AccessTokenScopeAdmin), string(AccessTokenScopeControlRw)},
			})
			require.NoError(t, err)

			expiredToken, err = querier.IssueAccessToken(ctx, db, IssueAccessTokenParams{
				Name:               "expired",
				TokenHash:          "expired-hash",
				NamespaceSelectors: []string{"ns-0", "ns-1"},
				Scopes:             []string{string(AccessTokenScopeAdmin), string(AccessTokenScopeControlRw)},
				ExpiresAt:          sql.NullTime{Valid: true, Time: time.Now().Add(-time.Hour)},
			})
			require.NoError(t, err)

			// This token will be revoked in the following subtest
			revokedToken, err = querier.IssueAccessToken(ctx, db, IssueAccessTokenParams{
				Name:               "revoked",
				TokenHash:          "revoked-hash",
				NamespaceSelectors: []string{"ns-0", "ns-1"},
				Scopes:             []string{string(AccessTokenScopeAdmin), string(AccessTokenScopeControlRw)},
			})
			require.NoError(t, err)
		})

		t.Run("RevokeAccessToken", func(t *testing.T) {
			err := querier.RevokeAccessToken(ctx, db, revokedToken.ID)
			require.NoError(t, err)
		})

		t.Run("GetAccessToken", func(t *testing.T) {
			t.Run("valid token", func(t *testing.T) {
				token, err := querier.GetAccessToken(ctx, db, validToken.ID)
				require.NoError(t, err)
				assert.Equal(t, validToken.ID, token.ID)
				assert.Equal(t, validToken.Name, token.Name)
				assert.Equal(t, validToken.NamespaceSelectors, token.NamespaceSelectors)
				assert.Equal(t, validToken.Scopes, token.Scopes)
				assert.Equal(t, validToken.IssuedAt, token.IssuedAt)
				assert.Equal(t, validToken.ExpiresAt, token.ExpiresAt)
			})
		})

		t.Run("ListAccessTokens", func(t *testing.T) {
			tokens, err := querier.ListAccessTokens(ctx, db, ListAccessTokensParams{
				IncludeExpired: false,
				IncludeRevoked: false,
			})
			require.NoError(t, err)
			assert.Equal(t, []ListAccessTokensRow{
				{
					ID:                 validToken.ID,
					Name:               validToken.Name,
					NamespaceSelectors: validToken.NamespaceSelectors,
					Scopes:             validToken.Scopes,
					IssuedAt:           validToken.IssuedAt,
					ExpiresAt:          validToken.ExpiresAt,
				},
				{
					ID:                 validNoExpiryToken.ID,
					Name:               validNoExpiryToken.Name,
					NamespaceSelectors: validNoExpiryToken.NamespaceSelectors,
					Scopes:             validNoExpiryToken.Scopes,
					IssuedAt:           validNoExpiryToken.IssuedAt,
					ExpiresAt:          validNoExpiryToken.ExpiresAt,
				},
			}, tokens)

			tokens, err = querier.ListAccessTokens(ctx, db, ListAccessTokensParams{IncludeExpired: true})
			require.NoError(t, err)
			assert.Len(t, tokens, 3)
			assert.Equal(t, expiredToken.ID, tokens[2].ID)
			assert.Equal(t, expiredToken.ExpiresAt, tokens[2].ExpiresAt)

			tokens, err = querier.ListAccessTokens(ctx, db, ListAccessTokensParams{IncludeRevoked: true})
			require.NoError(t, err)
			assert.Len(t, tokens, 3)
			assert.Equal(t, revokedToken.ID, tokens[2].ID)
			assert.False(t, tokens[2].RevokedAt.Time.After(time.Now()))
		})

		t.Run("AuthAccessToken", func(t *testing.T) {
			token, err := querier.AuthAccessToken(ctx, db, "valid-hash")
			require.NoError(t, err)
			assert.Equal(t, validToken.ID, token.ID)
			assert.Equal(t, validToken.Name, token.Name)
			assert.Equal(t, validToken.NamespaceSelectors, token.NamespaceSelectors)
			assert.Equal(t, validToken.Scopes, token.Scopes)

			token, err = querier.AuthAccessToken(ctx, db, "expired-hash")
			assert.Error(t, err)

			token, err = querier.AuthAccessToken(ctx, db, "revoked-hash")
			assert.Error(t, err)
		})
	})
}
