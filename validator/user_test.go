package validator

import (
	"testing"

	"github.com/clutterpot/clutterpot/model"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UserInput

func TestUserInputValidation(t *testing.T) {
	val := New()

	t.Run("Username", func(t *testing.T) { testUserInputValidation_Username(t, val) })
	t.Run("Email", func(t *testing.T) { testUserInputValidation_Email(t, val) })
	t.Run("Password", func(t *testing.T) { testUserInputValidation_Password(t, val) })
}

func testUserInputValidation_Username(t *testing.T, v *Validator) {
	t.Run("Required", func(t *testing.T) { testUserValidation_Username_required(t, v, model.UserInput{}) })
	t.Run("Username", func(t *testing.T) { testUserValidation_Username_username(t, v, model.UserInput{Username: "<test>"}) })
	t.Run("Min", func(t *testing.T) { testUserValidation_Username_min(t, v, model.UserInput{Username: randomString(3)}) })
	t.Run("Max", func(t *testing.T) {
		testUserValidation_Username_max(t, v, model.UserInput{Username: randomString(33)})
	})
	t.Run("Valid", func(t *testing.T) { testUserValidation_Username_valid(t, v, model.UserInput{Username: "test"}) })
}

func testUserInputValidation_Email(t *testing.T, v *Validator) {
	t.Run("Required", func(t *testing.T) { testUserValidation_Email_required(t, v, model.UserInput{}) })
	t.Run("Email", func(t *testing.T) { testUserValidation_Email_email(t, v, model.UserInput{Email: "test@example"}) })
	t.Run("Max", func(t *testing.T) {
		testUserValidation_Email_max(t, v, model.UserInput{Email: randomString(243) + "@example.com"})
	})
	t.Run("Valid", func(t *testing.T) { testUserValidation_Email_valid(t, v, model.User{Email: "test@example.com"}) })
}

func testUserInputValidation_Password(t *testing.T, v *Validator) {
	t.Run("Required", func(t *testing.T) { testUserValidation_Password_required(t, v, model.UserInput{}) })
	t.Run("Password", func(t *testing.T) {
		t.Run("ContainsLowercase", func(t *testing.T) {
			testUserValidation_Password_password_containsLowercase(t, v, model.UserInput{Password: "PASSWORD1!"})
		})
		t.Run("ContainsUppercase", func(t *testing.T) {
			testUserValidation_Password_password_containsUppercase(t, v, model.UserInput{Password: "password1!"})
		})
		t.Run("ContainsDigit", func(t *testing.T) {
			testUserValidation_Password_password_containsDigit(t, v, model.UserInput{Password: "Password!"})
		})
		t.Run("ContainsSpecial", func(t *testing.T) {
			testUserValidation_Password_password_containsSpecial(t, v, model.UserInput{Password: "Password1"})
		})
	})
	t.Run("Min", func(t *testing.T) { testUserValidation_Password_min(t, v, model.UserInput{Password: "Psswd1!"}) })
	t.Run("Max", func(t *testing.T) {
		testUserValidation_Password_max(t, v, model.UserInput{Password: randomString(252) + "Tt1!"})
	})
	t.Run("Valid", func(t *testing.T) { testUserValidation_Password_valid(t, v, model.UserInput{Password: "Password1!"}) })
}

// UserUpdateInput

func TestUserUpdateInputValidation(t *testing.T) {
	val := New()

	t.Run("Username", func(t *testing.T) { testUserUpdateInputValidation_Username(t, val) })
	t.Run("Email", func(t *testing.T) { testUserUpdateInputValidation_Email(t, val) })
	t.Run("Password", func(t *testing.T) { testUserUpdateInputValidation_Password(t, val) })
	t.Run("DisplayName", func(t *testing.T) { testUserUpdateInputValidation_DisplayName(t, val) })
	t.Run("Bio", func(t *testing.T) { testUserUpdateInputValidation_Bio(t, val) })
}

func testUserUpdateInputValidation_Username(t *testing.T, v *Validator) {
	var username string
	t.Run("Omitempty", func(t *testing.T) { testUserValidation_Username_omitempty(t, v, model.UserUpdateInput{}) })
	username = "<test>"
	t.Run("Username", func(t *testing.T) {
		testUserValidation_Username_username(t, v, model.UserUpdateInput{Username: &username})
	})
	username = randomString(3)
	t.Run("Min", func(t *testing.T) { testUserValidation_Username_min(t, v, model.UserUpdateInput{Username: &username}) })
	username = randomString(33)
	t.Run("Max", func(t *testing.T) { testUserValidation_Username_max(t, v, model.UserUpdateInput{Username: &username}) })
	username = "test"
	t.Run(("Valid"), func(t *testing.T) {
		testUserValidation_Username_valid(t, v, model.UserUpdateInput{Username: &username})
	})
}

func testUserUpdateInputValidation_Email(t *testing.T, v *Validator) {
	var email string
	t.Run("Omitempty", func(t *testing.T) { testUserValidation_Email_omitempty(t, v, model.UserUpdateInput{}) })
	email = "test@example"
	t.Run("Email", func(t *testing.T) { testUserValidation_Email_email(t, v, model.UserUpdateInput{Email: &email}) })
	email = randomString(243) + "@example.com"
	t.Run("Max", func(t *testing.T) { testUserValidation_Email_max(t, v, model.UserUpdateInput{Email: &email}) })
	email = "test@example.com"
	t.Run("Valid", func(t *testing.T) { testUserValidation_Email_valid(t, v, model.UserUpdateInput{Email: &email}) })
}

func testUserUpdateInputValidation_Password(t *testing.T, v *Validator) {
	var password string
	t.Run("Omitempty", func(t *testing.T) { testUserValidation_Password_omitempty(t, v, model.UserUpdateInput{}) })
	t.Run("Password", func(t *testing.T) {
		password = "PASSWORD1!"
		t.Run("ContainsLowercase", func(t *testing.T) {
			testUserValidation_Password_password_containsLowercase(t, v, model.UserUpdateInput{Password: &password})
		})
		password = "password1!"
		t.Run("ContainsUppercase", func(t *testing.T) {
			testUserValidation_Password_password_containsUppercase(t, v, model.UserUpdateInput{Password: &password})
		})
		password = "Password!"
		t.Run("ContainsDigit", func(t *testing.T) {
			testUserValidation_Password_password_containsDigit(t, v, model.UserUpdateInput{Password: &password})
		})
		password = "Password1"
		t.Run("ContainsSpecial", func(t *testing.T) {
			testUserValidation_Password_password_containsSpecial(t, v, model.UserUpdateInput{Password: &password})
		})
	})
	password = "Psswd1!"
	t.Run("Min", func(t *testing.T) { testUserValidation_Password_min(t, v, model.UserUpdateInput{Password: &password}) })
	password = randomString(252) + "Tt1!"
	t.Run("Max", func(t *testing.T) { testUserValidation_Password_max(t, v, model.UserUpdateInput{Password: &password}) })
	password = "Password1!"
	t.Run("Valid", func(t *testing.T) {
		testUserValidation_Password_valid(t, v, model.UserUpdateInput{Password: &password})
	})
}

func testUserUpdateInputValidation_DisplayName(t *testing.T, v *Validator) {
	var displayName string
	t.Run("Omitempty", func(t *testing.T) { testUserValidation_DisplayName_omitempty(t, v, model.UserUpdateInput{}) })
	displayName = " �\n"
	t.Run("DisplayName", func(t *testing.T) {
		testUserValidation_DisplayName_displayname(t, v, model.UserUpdateInput{DisplayName: &displayName})
	})
	displayName = randomString(33)
	t.Run("Max", func(t *testing.T) {
		testUserValidation_DisplayName_max(t, v, model.UserUpdateInput{DisplayName: &displayName})
	})
	displayName = "Test Display Name 1"
	t.Run("Valid", func(t *testing.T) {
		testUserValidation_DisplayName_valid(t, v, model.UserUpdateInput{DisplayName: &displayName})
	})
}

func testUserUpdateInputValidation_Bio(t *testing.T, v *Validator) {
	var bio string
	t.Run("Omitempty", func(t *testing.T) { testUserValidation_Bio_omitempty(t, v, model.UserUpdateInput{}) })
	bio = "\u00a0"
	t.Run("Printunicode", func(t *testing.T) { testUserValidation_Bio_printunicode(t, v, model.UserUpdateInput{Bio: &bio}) })
	bio = randomString(161)
	t.Run("Max", func(t *testing.T) { testUserValidation_Bio_max(t, v, model.UserUpdateInput{Bio: &bio}) })
	bio = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sed ornare augue, sed sodales nisi."
	t.Run("Valid", func(t *testing.T) { testUserValidation_Bio_valid(t, v, model.UserUpdateInput{Bio: &bio}) })
}

// Username

func testUserValidation_Username_required(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.Error(t, err, "validation should have failed on empty username")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"required\" tag")
	assert.Equal(t, "required", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"required\" tag")
}

func testUserValidation_Username_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.NoError(t, err, "validation should not have failed on empty username")
}

func testUserValidation_Username_username(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.Error(t, err, "validation should have failed on forbidden username characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"username\" tag")
	assert.Equal(t, "username", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"username\" tag")
}

func testUserValidation_Username_min(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.Error(t, err, "validation should have failed on too short username")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"min\" tag")
	assert.Equal(t, "min", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"min\" tag")
}

func testUserValidation_Username_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.Error(t, err, "validation should have failed on too long username")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testUserValidation_Username_valid(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Username")
	require.NoError(t, err, "validation should not have failed on valid username")
}

// Email

func testUserValidation_Email_required(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Email")
	require.Error(t, err, "validation should have failed on empty email")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"required\" tag")
	assert.Equal(t, "required", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"required\" tag")
}

func testUserValidation_Email_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Email")
	require.NoError(t, err, "validation should not have failed on empty email")

}

func testUserValidation_Email_email(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Email")
	require.Error(t, err, "validation should have failed on invalid email address")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"email\" tag")
	assert.Equal(t, "email", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"email\" tag")
}

func testUserValidation_Email_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Email")
	require.Error(t, err, "validation should have failed on too long email")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testUserValidation_Email_valid(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Email")
	require.NoError(t, err, "validation should not have failed on valid email")
}

// Password

func testUserValidation_Password_required(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on empty password")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"required\" tag")
	assert.Equal(t, "required", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"required\" tag")
}

func testUserValidation_Password_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.NoError(t, err, "validation should not have failed on empty password")
}

func testUserValidation_Password_password_containsLowercase(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on password with no lowercase character")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"password\" tag")
	assert.Equal(t, "password", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"password\" tag")
}

func testUserValidation_Password_password_containsUppercase(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on password with no uppercase characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"password\" tag")
	assert.Equal(t, "password", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"password\" tag")
}

func testUserValidation_Password_password_containsDigit(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on password with no digit")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"password\" tag")
	assert.Equal(t, "password", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"password\" tag")
}

func testUserValidation_Password_password_containsSpecial(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on password with no special character")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"password\" tag")
	assert.Equal(t, "password", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"password\" tag")
}

func testUserValidation_Password_min(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on too short password")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"min\" tag")
	assert.Equal(t, "min", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"min\" tag")
}

func testUserValidation_Password_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.Error(t, err, "validation should have failed on too long password")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testUserValidation_Password_valid(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Password")
	require.NoError(t, err, "validation should not have failed on valid password")
}

// DisplayName

func testUserValidation_DisplayName_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "DisplayName")
	require.NoError(t, err, "validation should not have failed on empty display name")
}

func testUserValidation_DisplayName_displayname(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "DisplayName")
	require.Error(t, err, "validation should have failed on forbidden display name characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"displayname\" tag")
	assert.Equal(t, "displayname", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"displayname\" tag")
}

func testUserValidation_DisplayName_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "DisplayName")
	require.Error(t, err, "validation should have failed on too long display name")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testUserValidation_DisplayName_valid(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "DisplayName")
	require.NoError(t, err, "validation should not have failed on valid display name")
}

// Bio

func testUserValidation_Bio_omitempty(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Bio")
	require.NoError(t, err, "validation should not have failed on empty bio")
}

func testUserValidation_Bio_printunicode(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Bio")
	require.Error(t, err, "validation should have failed on bio with non printable unicode characters")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"printunicode\" tag")
	assert.Equal(t, "printunicode", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"printunicode\" tag")
}

func testUserValidation_Bio_max(t *testing.T, v *Validator, m any) {
	err := v.val.StructPartial(m, "Bio")
	require.Error(t, err, "validation should have failed on too long bio")
	assert.Equal(t, len(err.(validator.ValidationErrors)), 1, "validation should have failed only on \"max\" tag")
	assert.Equal(t, "max", err.(validator.ValidationErrors)[0].Tag(), "validation should have failed on \"max\" tag")
}

func testUserValidation_Bio_valid(t *testing.T, val *Validator, m any) {
	err := val.val.StructPartial(m, "Bio")
	require.NoError(t, err, "validation should not have failed on valid bio")
}
