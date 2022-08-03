package validator

import (
	"testing"

	"github.com/clutterpot/clutterpot/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileInputValidation(t *testing.T) {
	val := New()

	t.Run("Filename", func(t *testing.T) { testFileInputValidation_Filename(t, val) })
}

func testFileInputValidation_Filename(t *testing.T, v *Validator) {
	t.Run("Required", func(t *testing.T) { testFileValidation_Filename_required(t, v, model.FileInput{}) })
	t.Run("Filename", func(t *testing.T) { testFileValidation_Filename_filename(t, v, model.FileInput{Filename: "/>|:&"}) })
	t.Run("Printunicode", func(t *testing.T) {
		testFileValidation_Filename_printunicode(t, v, model.FileInput{Filename: "\u00a0"})
	})
	t.Run("Max", func(t *testing.T) {
		testFileValidation_Filename_max(t, v, model.FileInput{Filename: randomString(256)})
	})
	t.Run("Valid", func(t *testing.T) { testFileValidation_Filename_valid(t, v, model.FileInput{Filename: "test"}) })
}

func TestFileUpdateInputValidation(t *testing.T) {
	val := New()

	t.Run("Filename", func(t *testing.T) { testFileUpdateInputValidation_Filename(t, val) })
}

func testFileUpdateInputValidation_Filename(t *testing.T, v *Validator) {
	var filename string
	t.Run("Omitempty", func(t *testing.T) { testFileValidation_Filename_omitempty(t, v, model.FileUpdateInput{}) })
	filename = "/>|:&\n"
	t.Run("Filename", func(t *testing.T) {
		testFileValidation_Filename_filename(t, v, model.FileUpdateInput{Filename: &filename})
	})
	filename = "\u00a0"
	t.Run("Printunicode", func(t *testing.T) {
		testFileValidation_Filename_printunicode(t, v, model.FileUpdateInput{Filename: &filename})
	})
	filename = ""
	t.Run("Min", func(t *testing.T) { testFileValidation_Filename_min(t, v, model.FileUpdateInput{Filename: &filename}) })
	filename = randomString(256)
	t.Run("Max", func(t *testing.T) { testFileValidation_Filename_max(t, v, model.FileUpdateInput{Filename: &filename}) })
	filename = "test"
	t.Run("Valid", func(t *testing.T) {
		testFileValidation_Filename_valid(t, v, model.FileUpdateInput{Filename: &filename})
	})
}

// Filename

func testFileValidation_Filename_required(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.Error(t, err, "validation should have failed on empty filename")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"required\" tag")
	assert.Equal(t, "required", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"required\" tag")
}

func testFileValidation_Filename_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.NoError(t, err, "validation should not have failed on empty filename")
}

func testFileValidation_Filename_filename(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.Error(t, err, "validation should have failed on forbidden filename characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"filename\" tag")
	assert.Equal(t, "filename", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"filename\" tag")
}

func testFileValidation_Filename_printunicode(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.Error(t, err, "validation should have failed on filename with non printable unicode characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"printunicode\" tag")
	assert.Equal(t, "printunicode", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"printunicode\" tag")
}

func testFileValidation_Filename_min(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.Error(t, err, "validation should have failed on too short filename")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"min\" tag")
	assert.Equal(t, "min", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"min\" tag")
}

func testFileValidation_Filename_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.Error(t, err, "validation should have failed on too long filename")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testFileValidation_Filename_valid(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Filename")
	require.NoError(t, err, "validation should not have failed on valid filename")
}
