// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj-up/pkg/common"
)

func Test_Undo(t *testing.T) {
	testCmd := initCmd()
	testCmd.SetArgs([]string{})
	err := testCmd.Execute()
	require.NoError(t, err)
	testCmd = setEnvCmd()
	testCmd.SetArgs([]string{"authservice", "test1=1"})
	err = testCmd.Execute()
	require.NoError(t, err)
	newTemplateBytes, err := common.Store.RestoreLatestVersion()
	require.False(t, bytes.Contains(newTemplateBytes, []byte("test1")))
	require.Nil(t, err)
	cleanup(t)

	testCmd = initCmd()
	testCmd.SetArgs([]string{})
	err = testCmd.Execute()
	require.NoError(t, err)
	testCmd = setEnvCmd()
	testCmd.SetArgs([]string{"authservice", "test1=1"})
	err = testCmd.Execute()
	require.NoError(t, err)
	testCmd.SetArgs([]string{"authservice", "test2=2"})
	err = testCmd.Execute()
	require.NoError(t, err)
	newTemplateBytes, err = common.Store.RestoreLatestVersion()
	require.True(t, bytes.Contains(newTemplateBytes, []byte("test1")))
	require.Nil(t, err)
	cleanup(t)

	testCmd = initCmd()
	testCmd.SetArgs([]string{})
	err = testCmd.Execute()
	require.NoError(t, err)
	testCmd = setEnvCmd()
	testCmd.SetArgs([]string{"authservice", "test1=1"})
	err = testCmd.Execute()
	require.NoError(t, err)
	testCmd.SetArgs([]string{"authservice", "test2=2"})
	err = testCmd.Execute()
	require.NoError(t, err)
	testCmd.SetArgs([]string{"authservice", "test3=3"})
	err = testCmd.Execute()
	require.NoError(t, err)
	newTemplateBytes, err = common.Store.RestoreLatestVersion()
	require.True(t, bytes.Contains(newTemplateBytes, []byte("test2")))
	require.Nil(t, err)
	testCmd = undoCmd()
	err = testCmd.Execute()
	require.NoError(t, err)
	newTemplateBytes, err = common.Store.RestoreLatestVersion()
	require.False(t, bytes.Contains(newTemplateBytes, []byte("test2")))
	require.Nil(t, err)
	testCmd = undoCmd()
	require.NoError(t, err)
	err = testCmd.Execute()
	require.NoError(t, err)
	newTemplateBytes, err = common.Store.RestoreLatestVersion()
	require.Nil(t, newTemplateBytes)
	require.Error(t, err) // DB is empty, no history to restore
	cleanup(t)
}

func cleanup(t *testing.T) {
	err := os.Remove("docker-compose.yaml")
	require.NoError(t, err)
	path, _ := filepath.Abs("./database")
	err = os.RemoveAll(path)
	require.NoError(t, err)
}
