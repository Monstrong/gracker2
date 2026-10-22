package service

// DeleteUser TODO: проверка: выполнить только если id == callerID || callerRole == admin
// password validate: if err == bcrypt.ErrPasswordTooLong (>72byte), 