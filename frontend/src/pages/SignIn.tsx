import { useState, FormEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import { signIn } from "../api/auth";
import { ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import styles from "./SignIn.module.css";

interface FormErrors {
  username?: string;
  password?: string;
}

export default function SignIn() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<FormErrors>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function validate(): FormErrors {
    const next: FormErrors = {};
    if (!username.trim()) next.username = "Введите username";
    if (!password) next.password = "Введите пароль";
    return next;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setServerError(null);

    const validationErrors = validate();
    setErrors(validationErrors);
    if (Object.keys(validationErrors).length > 0) return;

    setLoading(true);
    try {
      const { token } = await signIn({ username, password });
      login(token);
      navigate("/");
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        setServerError("Неверный username или пароль");
      } else {
        setServerError("Что-то пошло не так, попробуйте позже");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className={styles.wrapper}>
      <form className={styles.card} onSubmit={handleSubmit} noValidate>
        <h1 className={styles.title}>Вход</h1>

        <label className={styles.field}>
          <span>Username</span>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className={errors.username ? styles.inputError : undefined}
          />
          {errors.username && <span className={styles.error}>{errors.username}</span>}
        </label>

        <label className={styles.field}>
          <span>Пароль</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={errors.password ? styles.inputError : undefined}
          />
          {errors.password && <span className={styles.error}>{errors.password}</span>}
        </label>

        {serverError && <div className={styles.serverError}>{serverError}</div>}

        <button type="submit" className={styles.submit} disabled={loading}>
          {loading ? "Входим..." : "Войти"}
        </button>

        <p className={styles.switch}>
          Нет аккаунта? <Link to="/sign-up">Зарегистрироваться</Link>
        </p>
      </form>
    </div>
  );
}