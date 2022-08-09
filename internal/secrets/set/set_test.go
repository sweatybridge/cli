package set

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/supabase/cli/internal/secrets/list"
	"github.com/supabase/cli/internal/testing/apitest"
	"github.com/supabase/cli/internal/utils"
	"gopkg.in/h2non/gock.v1"
)

func TestSecretSetCommand(t *testing.T) {
	dummy := list.Secret{Name: "my-name", Value: "my-value"}
	dummyEnv := dummy.Name + "=" + dummy.Value

	t.Run("Sets secret via cli args", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Setup valid project ref
		project := apitest.RandomProjectRef()
		err = afero.WriteFile(fsys, utils.ProjectRefPath, []byte(project), 0644)
		require.NoError(t, err)
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.Off()
		gock.New("https://api.supabase.io").
			Post("/v1/projects/" + project + "/secrets").
			MatchType("json").
			JSON([]list.Secret{dummy}).
			Reply(200)
		// Run test
		assert.NoError(t, Run("", []string{dummyEnv}, fsys))
	})

	t.Run("Sets secret value via env file", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Setup valid project ref
		project := apitest.RandomProjectRef()
		err = afero.WriteFile(fsys, utils.ProjectRefPath, []byte(project), 0644)
		require.NoError(t, err)
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Setup dotenv file
		tmpfile, err := os.CreateTemp("", "secret")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		_, err = tmpfile.Write([]byte(dummyEnv))
		require.NoError(t, err)
		// Flush pending mocks after test execution
		defer gock.Off()
		gock.New("https://api.supabase.io").
			Post("/v1/projects/" + project + "/secrets").
			MatchType("json").
			JSON([]list.Secret{dummy}).
			Reply(200)
		// Run test
		assert.NoError(t, Run(tmpfile.Name(), []string{}, fsys))
	})

	t.Run("throws error on missing config file", func(t *testing.T) {
		assert.Error(t, Run("", []string{}, afero.NewMemMapFs()))
	})

	t.Run("throws error on missing project ref", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Run test
		assert.Error(t, Run("", []string{}, fsys))
	})

	t.Run("throws error on missing access token", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		_, err = fsys.Create(utils.ProjectRefPath)
		require.NoError(t, err)
		// Run test
		assert.Error(t, Run("", []string{}, fsys))
	})

	t.Run("throws error on empty secret", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Setup valid project ref
		project := apitest.RandomProjectRef()
		err = afero.WriteFile(fsys, utils.ProjectRefPath, []byte(project), 0644)
		require.NoError(t, err)
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Run test
		assert.Error(t, Run("", []string{}, fsys))
	})

	t.Run("throws error on network error", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Setup valid project ref
		project := apitest.RandomProjectRef()
		err = afero.WriteFile(fsys, utils.ProjectRefPath, []byte(project), 0644)
		require.NoError(t, err)
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.Off()
		gock.New("https://api.supabase.io").
			Post("/v1/projects/" + project + "/secrets").
			MatchType("json").
			JSON([]list.Secret{dummy}).
			ReplyError(errors.New("network error"))
		// Run test
		assert.Error(t, Run("", []string{dummyEnv}, fsys))
	})

	t.Run("throws error on server unavailable", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create(utils.ConfigPath)
		require.NoError(t, err)
		// Setup valid project ref
		project := apitest.RandomProjectRef()
		err = afero.WriteFile(fsys, utils.ProjectRefPath, []byte(project), 0644)
		require.NoError(t, err)
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.Off()
		gock.New("https://api.supabase.io").
			Post("/v1/projects/" + project + "/secrets").
			MatchType("json").
			JSON([]list.Secret{dummy}).
			Reply(500).
			JSON(map[string]string{"message": "unavailable"})
		// Run test
		assert.Error(t, Run("", []string{dummyEnv}, fsys))
	})
}