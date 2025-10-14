CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL
);
-- Password: 'awnoidroppeditinthewater' (bcrypt hash)

INSERT INTO users (username, password_hash) VALUES (
  'user-admin', 
  '$2a$12$S7dn3Y/ltqTn6kNwhi0Hee/hvcgC9Z4jFiWhU30J19uQnl48BsURO'
);

-- Password: 'iamamovinggroovingjammingsinginggummybear' (bcrypt hash)

INSERT INTO users (username, password_hash) VALUES (
  'server-admin', 
  '$2a$12$xapBUEBgT0FFgGFt7M.yveuPOHTxSszLSgo9xcV9VKExQRcNNfXxS'
);
