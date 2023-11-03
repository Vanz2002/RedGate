interface LoginResponse {
  sessionID: string;
  accessToken: string;
  accessTokenEx: Date;
  refreshToken: string;
  refreshTokenEx: Date;
  userID: string;
  username: string;
}

interface LoginResponseJson {
  session_id: string;
  access_token: string;
  access_token_expire: string;
  refresh_token: string;
  refresh_token_expire: string;
  user_id: string;
  username: string;
}

function toLoginResponse(json: LoginResponseJson): LoginResponse {
  return {
    sessionID: json.session_id,
    accessToken: json.access_token,
    accessTokenEx: new Date(json.access_token_expire),
    refreshToken: json.refresh_token,
    refreshTokenEx: new Date(json.refresh_token_expire),
    userID: json.user_id,
    username: json.username,
  };
}

interface RenewToken {
  accessTokenExpire: Date;
  refreshToken: string;
}