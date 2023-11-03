import React, { useEffect, useState } from "react";
import "../../App.css";
import "../../index.css";
import axios from "axios";
import { useCookies } from "react-cookie";
import checkAndRenewToken from "./renewToken";

const LoginForm: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSignUp, setIsSignUp] = useState(false);
  const [username, setUsername] = useState("");
  const [cookies, setCookie] = useCookies([
    "accessToken",
    "refreshToken",
    "accessTokenEx",
  ]);

  const [success, setSuccess] = useState(false);
  const [accounts, setAccounts] = useState<Account[]>([]);

  const instance = axios.create({
    baseURL: "http://127.0.0.1:4444/", // Replace with your API base UR
  });

  useEffect(() => {
    const renewTokenParam: RenewToken = {
      accessTokenExpire: new Date(cookies.accessTokenEx),
      refreshToken: cookies.refreshToken,
    };
    // Call the checkAndRenewToken function before making any API call
    checkAndRenewToken(renewTokenParam, setCookie).then(() => {
      instance
        .get("/account/list", {
          headers: {
            Authorization: `Bearer ${cookies.accessToken}`,
          },
        })
        .then((response) => {
          if (response.status == 200) {
            setAccounts(response.data); // Set the data in state
          }
          console.log(response.data);
        })
        .catch((error) => console.error("Error fetching data:", error));
    });
  }, [success]);

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
  };

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    const formData = new URLSearchParams();
    formData.append("username", username);
    formData.append("email", email);
    formData.append("password", password);

    try {
      await instance.post("http://127.0.0.1:4444/auth/signup", formData, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      });
    } catch (error) {
      console.error("Form submission error:", error);
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    const formData = new URLSearchParams();
    formData.append("email", email);
    formData.append("password", password);

    try {
      const resp = await instance.post(
        "http://127.0.0.1:4444/auth/login",
        formData,
        {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
        }
      );

      const { access_token, access_token_expire, refresh_token } =
        resp.data as LoginResponseJson;
      setCookie("accessToken", access_token);
      setCookie("accessTokenEx", access_token_expire);
      setCookie("refreshToken", refresh_token);
      setSuccess(true);
    } catch (error) {
      console.error("Form submission error:", error);
    }
  };

  return (
    <div className="max-w-md mx-auto my-10 p-6 bg-white rounded-md shadow-lg">
      <h2 className="text-[100px] font-bold mb-6 text-green-700">RED-GATE</h2>
      <form onSubmit={isSignUp ? handleSignUp : handleLogin}>
        <div className="mb-4">
          <label
            htmlFor="email"
            className="block text-gray-700 text-sm font-bold mb-2"
          >
            Email
          </label>
          <input
            type="email"
            id="email"
            className="w-full p-2 border rounded"
            placeholder="Enter your email"
            value={email}
            onChange={handleEmailChange}
            required
          />
        </div>
        <div className="mb-4">
          <label
            htmlFor="password"
            className="block text-gray-700 text-sm font-bold mb-2"
          >
            Password
          </label>
          <input
            type="password"
            id="password"
            className="w-full p-2 border rounded"
            placeholder="Enter your password"
            value={password}
            onChange={handlePasswordChange}
            required
          />
        </div>
        {isSignUp && (
          <div className="mb-4">
            <label
              htmlFor="username"
              className="block text-gray-700 text-sm font-bold mb-2"
            >
              Username
            </label>
            <input
              type="text"
              id="username"
              className="w-full p-2 border rounded"
              placeholder="Enter your username"
              value={username}
              onChange={handleUsernameChange}
              required
            />
          </div>
        )}
        <button
          type="submit"
          className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          {isSignUp ? "Sign Up" : "Login"}
        </button>
        <div className="mt-4">
          <button
            type="button"
            className="text-gray-700 hover:underline"
            onClick={() => setIsSignUp((prev) => !prev)}
          >
            {isSignUp
              ? "Already have an account? Login"
              : "Don't have an account? Sign Up"}
          </button>
        </div>
      </form>
      <div>
        <h2>Fetched Data:</h2>
        <ul>
          {accounts.map((account, index) => (
            <li key={index}>{account.account_id}</li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default LoginForm;
