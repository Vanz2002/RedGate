import {
  BrowserRouter,
  Navigate,
  Route,
  Router,
  RouterProvider,
  Routes,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";
import "./App.css";
import "./index.css";
import LoginForm from "./page/login";
import HomePage from "./page/home_page";
import { AuthProvider } from "./auth/AuthProvider";

function App() {
  return (
      <AuthProvider>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginForm />} />
        </Routes>
      </AuthProvider>
  );
}

export default App;
